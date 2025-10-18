package main

import (
	"fmt"
	"github.com/rookgm/gophkeeper/config"
	"github.com/rookgm/gophkeeper/internal/build"
	"github.com/rookgm/gophkeeper/internal/client/cli"
)

// application build info
var (
	// BuildVersion is application build version
	BuildVersion = "N/A"
	// BuildDate is application build date
	BuildDate = "N/A"
	// BuildCommit is application build commit
	BuildCommit = "N/A"
)

func main() {

	// load client config
	cfg, err := config.NewClientConfig()
	if err != nil {
		panic(err)
	}

	// TODO
	fmt.Println(cfg)

	info := build.NewBuildInfo(BuildVersion, BuildDate, BuildCommit)
	cli.Execute(info)
}
