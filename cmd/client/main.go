package main

import (
	"fmt"
	"github.com/rookgm/gophkeeper/config"
	"github.com/rookgm/gophkeeper/internal/build"
	"github.com/rookgm/gophkeeper/internal/client/api"
	"github.com/rookgm/gophkeeper/internal/client/cli"
	"github.com/rookgm/gophkeeper/internal/client/crypto"
	"github.com/rookgm/gophkeeper/internal/client/service"
	"log"
	"os"
	"path"
)

const tokenFileName = "token"

// application build info
var (
	// BuildVersion is application build version
	BuildVersion = "N/A"
	// BuildDate is application build date
	BuildDate = "N/A"
	// BuildCommit is application build commit
	BuildCommit = "N/A"
)

// createConfigDir creates configuration dir
// if path is empty, the creates dir at default root directory
// to use for user-specific configuration data.
func createConfigDir(dir string) (string, error) {
	root := dir
	if root == "" {
		cfg, _ := os.UserConfigDir()
		root = path.Join(cfg, "gophkeeper")
	}

	// create client configuration dir
	if err := os.MkdirAll(root, 0700); err != nil {
		log.Fatalf("error creating config dir %s: %v", root, err)
	}
	return root, nil
}

func main() {
	// load client config
	cfg, err := config.NewClientConfig()
	if err != nil {
		log.Fatalf("error creating new client config: %v", err)
	}

	cfgPath, err := createConfigDir(cfg.ConfigDir)
	if err != nil {
		log.Fatalf("error creating config dir: %v", err)
	}

	// dependency injection
	// application build info
	buildInfo := build.NewBuildInfo(BuildVersion, BuildDate, BuildCommit)
	// create AES encryptor
	encryptor := crypto.NewAESEncryptor()
	// create client based on http
	apiClient := api.NewClient(cfg.ServerAddress)
	// create token service
	tokener := service.NewTokenService(path.Join(cfgPath, tokenFileName))
	// create user service
	userService := service.NewUserService(apiClient, tokener)
	// create secret service
	secretService := service.NewSecretService(apiClient, encryptor, tokener)
	// create client CLI
	clientCli := cli.NewRootCmd(userService, secretService, buildInfo)

	// run CLI
	if err := clientCli.Execute(); err != nil {
		fmt.Printf("error running client cli %v\n", err)
		os.Exit(1)
	}
}
