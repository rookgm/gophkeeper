package main

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rookgm/gophkeeper/internal/server/auth"
	"github.com/rookgm/gophkeeper/internal/server/http/handler"
	"github.com/rookgm/gophkeeper/internal/server/http/middleware"
	"github.com/rookgm/gophkeeper/internal/server/repository"
	"github.com/rookgm/gophkeeper/internal/server/repository/postgres"
	"github.com/rookgm/gophkeeper/internal/server/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rookgm/gophkeeper/config"
	"github.com/rookgm/gophkeeper/internal/logger"
	"go.uber.org/zap"
)

const (
	serverCertFileName = "cert/server.crt"
	serverKeyFileName  = "cert/server.key"
)

const shutdownTimeout = 5 * time.Second

func main() {
	// initialize config
	cfg, err := config.Initialize()
	if err != nil {
		panic("Error initialize config")
	}

	// initialize logger
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic("Error initialize logger")
	}

	// create context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// check auth token key
	if cfg.AuthTokenKey == "" {
		logger.Log.Fatal("Error token key is empty", zap.Error(err))
	}

	// decode authentification token in hex string
	tokenKey, err := hex.DecodeString(cfg.AuthTokenKey)
	if err != nil {
		logger.Log.Fatal("Error extracting token key", zap.Error(err))
	}

	// initialize database
	db, err := postgres.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		logger.Log.Fatal("Error initializing database", zap.Error(err))
	}
	defer db.Close()

	// migrate database
	if err := db.Migrate(); err != nil {
		logger.Log.Fatal("Error migrating database", zap.Error(err))
	}

	// dependency injection
	// token
	token := auth.NewAuthToken(tokenKey)
	// user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, token)

	// auth
	authService := service.NewAuthService(userRepo, token)
	authHandler := handler.NewAuthHandler(authService)

	// secret
	secretRepo := repository.NewSecretRepository(db)
	secretService := service.NewSecretService(secretRepo)
	secretHandler := handler.NewSecretHandler(secretService)

	// routes
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.Logging(logger.Log))
	router.Use(middleware.Gzip)

	// register user
	router.Post("/api/user/register", userHandler.RegisterUser())
	// login user
	router.Post("/api/user/login", authHandler.LoginUser())

	// routes that require authentication
	router.Group(func(group chi.Router) {
		group.Use(middleware.Auth(token))
		group.Post("/api/user/secrets", secretHandler.CreateUserSecret)
		group.Get("/api/user/secrets/{id}", secretHandler.GetUserSecret)
		group.Put("/api/user/secrets/{id}", secretHandler.UpdateUserSecret)
		group.Delete("/api/user/secrets/{id}", secretHandler.DeleteUserSecret)
	})

	// set server parameters
	srv := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// check existing server's key files
	if _, err := os.Stat(serverCertFileName); errors.Is(err, os.ErrNotExist) {
		logger.Log.Fatal("server cert file is not exist", zap.Error(err))
	}
	if _, err := os.Stat(serverKeyFileName); errors.Is(err, os.ErrNotExist) {
		logger.Log.Fatal("server key file is not exist", zap.Error(err))
	}

	logger.Log.Info("Starting server...")
	go func() {
		logger.Log.Info("Starting HTTPs server", zap.String("address", cfg.Address))
		if err := http.ListenAndServeTLS(cfg.Address, serverCertFileName, serverKeyFileName, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal("Error starting https server", zap.Error(err))
		}
	}()

	logger.Log.Info("Server is started", zap.String("addr", cfg.Address))

	<-ctx.Done()

	logger.Log.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// shutdown server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Error shutdown server", zap.Error(err))
	}

	// close db
	if db != nil {
		db.Close()
	}

	logger.Log.Info("server is finished")
}
