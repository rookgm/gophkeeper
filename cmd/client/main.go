package main

import (
	"github.com/rookgm/gophkeeper/config"
	"github.com/rookgm/gophkeeper/internal/client/cmd"
	"github.com/rookgm/gophkeeper/internal/logger"
	"go.uber.org/zap"
)

func main() {

	// load client config
	cfg, err := config.NewClientConfig()
	if err != nil {
		panic(err)
	}

	// initialize logger
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic("Error initialize logger")
	}

	logger.Log.Info("Starting client", zap.String("address", cfg.ServerAddress))

	cmd.Execute()
}
