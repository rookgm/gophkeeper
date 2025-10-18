package config

import (
	"flag"
	"os"
	"sync"
)

// client default values
const (
	defaultServerAddress = ":8443"
	defaultLogLevel      = "debug"
	defaultStoragePath   = "./gophkeeper/storage"
)

// ClientConfig contains client configuration information
type ClientConfig struct {
	ServerAddress string
	LogLevel      string
	StoragePath   string
}

// singleton
var (
	once      sync.Once
	singleton *ClientConfig
)

func NewClientConfig() (*ClientConfig, error) {
	once.Do(func() {
		cfg := ClientConfig{}

		// init flags
		flag.StringVar(&cfg.ServerAddress, "a", defaultServerAddress, "gophkeeper server address")
		flag.StringVar(&cfg.LogLevel, "l", defaultLogLevel, "gophkeeper client log level")
		flag.StringVar(&cfg.StoragePath, "f", defaultStoragePath, "gophkeeper client storage path")

		flag.Parse()

		// get config from envs
		// sets server address
		if serverAddrEnv := os.Getenv("GOPHKEEPER_SERVER_ADDRESS"); serverAddrEnv != "" {
			cfg.ServerAddress = serverAddrEnv
		}
		// sets client logging level
		if logLevelEnv := os.Getenv("CLIENT_LOG_LEVEL"); logLevelEnv != "" {
			cfg.LogLevel = logLevelEnv
		}
		// sets client storage path
		if storagePathEnv := os.Getenv("CLIENT_STORAGE_PATH"); storagePathEnv != "" {
			cfg.StoragePath = storagePathEnv
		}

		singleton = &cfg
	})

	return singleton, nil
}
