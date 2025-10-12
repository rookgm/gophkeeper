package main

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rookgm/gophkeeper/internal/auth"
	"github.com/rookgm/gophkeeper/internal/http/handler"
	"github.com/rookgm/gophkeeper/internal/http/middleware"
	"github.com/rookgm/gophkeeper/internal/repository"
	"github.com/rookgm/gophkeeper/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rookgm/gophkeeper/config"
	"github.com/rookgm/gophkeeper/internal/logger"
	"github.com/rookgm/gophkeeper/internal/repository/postgres"
	"go.uber.org/zap"
)

const authTokenKey = "f53ac685bbceebd75043e6be2e06ee07"

const (
	serverCertFileName = "cert/server.crt"
	serverKeyFileName  = "cert/server.key"
)

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

	// decode authentification token key
	tokenKey, err := hex.DecodeString(authTokenKey)
	if err != nil {
		logger.Log.Fatal("Error extracting token key", zap.Error(err))
	}
	token := auth.NewAuthToken(tokenKey)

	// dependency injection
	// user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, token)

	// auth
	authService := service.NewAuthService(userRepo, token)
	authHandler := handler.NewAuthHandler(authService)

	// routes
	router := chi.NewRouter()

	router.Use(middleware.Logging(logger.Log))
	// register user
	router.Post("/api/user/register", userHandler.RegisterUser())
	// login user
	router.Post("/api/user/login", authHandler.LoginUser())

	// routes that require authentication
	router.Group(func(group chi.Router) {
		group.Use(middleware.Auth(token))
	})

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

	// TODO add graceful shutdown
}
