package config

import (
	"encoding/json"
	"flag"
	"os"
)

const (
	defaultServerAddr     = ":8443"
	defaultServerLogLevel = "debug"
)

// ServerConfig contains server configurations
type ServerConfig struct {
	// Address is server address(host:port)
	Address string
	// DatabaseDSN is data source name(postgres://user:pass@tcp(localhost:5555)/dbname)
	DatabaseDSN string
	// LogLevel is logging level
	LogLevel string
	// ConfigPath is configuration file path
	ConfigPath string
}

// ServerOption is server config func option
type ServerOption func(*ServerConfig)

// WithServerAddr sets server address in Config
func WithServerAddr(addr string) ServerOption {
	return func(c *ServerConfig) {
		if addr != "" {
			c.Address = addr
		}
	}
}

// WithLogLevel sets server logging level
func WithLogLevel(level string) ServerOption {
	return func(c *ServerConfig) {
		if level != "" {
			c.LogLevel = level
		}
	}
}

// WithDatabaseDSN sets server data source name
func WithDatabaseDSN(dsn string) ServerOption {
	return func(c *ServerConfig) {
		if dsn != "" {
			c.DatabaseDSN = dsn
		}
	}
}

// serverConfigJSON presents server config in json format
type serverConfigJSON struct {
	ServerAddress string `json:"server_address"`
	DatabaseDSN   string `json:"database_dsn"`
}

// FromFile loads server config from file in JSON format
func FromFile(name string) ServerOption {
	return func(c *ServerConfig) {
		if name == "" {
			return
		}

		b, err := os.ReadFile(name)
		if err != nil {
			return
		}

		cfg := serverConfigJSON{}

		err = json.Unmarshal(b, &cfg)
		if err != nil {
			return
		}

		WithServerAddr(cfg.ServerAddress)(c)
		WithDatabaseDSN(cfg.DatabaseDSN)(c)
	}
}

func FromEnv() ServerOption {
	return func(c *ServerConfig) {
		// sets server address
		if serverAddrEnv := os.Getenv("SERVER_ADDRESS"); serverAddrEnv != "" {
			WithServerAddr(serverAddrEnv)(c)
		}
		// sets database source name
		if dataBaseDSNEnv := os.Getenv("SERVER_DATABASE_DSN"); dataBaseDSNEnv != "" {
			WithDatabaseDSN(dataBaseDSNEnv)(c)
		}
		// sets logging level
		if logLevelEnv := os.Getenv("SERVER_LOG_LEVEL"); logLevelEnv != "" {
			WithLogLevel(logLevelEnv)(c)
		}
	}
}

// FromCommandLine gets server configuration from command line
func FromCommandLine(args *ServerConfig) ServerOption {
	return func(c *ServerConfig) {
		WithServerAddr(args.Address)(c)
		WithDatabaseDSN(args.DatabaseDSN)(c)
		WithLogLevel(args.LogLevel)(c)
	}
}

// parseCommandLine parses command line arguments
func parseCommandLine(cfg *ServerConfig) {
	flag.StringVar(&cfg.Address, "a", "", "server address")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "database address")
	flag.StringVar(&cfg.LogLevel, "l", "", "log level")
	flag.StringVar(&cfg.ConfigPath, "c", "", "load config from file")

	flag.Parse()
}

// New returns new server Config. It parses command line, environment variables and file.
func New(opts ...ServerOption) (*ServerConfig, error) {
	// set defaults values
	cfg := &ServerConfig{
		Address:  defaultServerAddr,
		LogLevel: defaultServerLogLevel,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg, nil
}

// Initialize initializes the server configuration
func Initialize() (*ServerConfig, error) {
	args := &ServerConfig{}
	// parse server command line
	parseCommandLine(args)
	return New(
		// low priority
		FromFile(args.ConfigPath),
		// medium priority
		FromEnv(),
		// height priority
		FromCommandLine(args),
	)
}
