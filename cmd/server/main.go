package main

import (
	"github.com/rookgm/gophkeeper/config"
	"github.com/rookgm/gophkeeper/internal/logger"
)

func main() {
	// initialize config
	cfg, err := config.Initialize()
	if err != nil {
		panic("Can not initialize config")
	}

	// initialize logger
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic("Can not initialize logger")
	}

	logger.Log.Info("info")

	// initialize database
	// -d postgres://gophermart:gophermart@postgres:5432/gophermart?sslmode=disable
	// -d postgres://gophkeeper:gophkeeper@postgres:5432/gophkeeper?sslmode=disable
}
