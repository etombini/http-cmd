package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/etombini/http-cmd/pkg/config"
	"github.com/etombini/http-cmd/pkg/server"
	"github.com/etombini/http-cmd/pkg/version"
)

func main() {

	versionFlag := flag.Bool("version", false, "Get version")
	configFlag := flag.String("config", config.DefaultConfPath, "Configuration file ["+config.DefaultConfPath+"]")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version:\t%s\n", version.Version())
		fmt.Printf("Build:  \t%s\n", version.Build())
		return
	}

	cfg, err := config.New(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	server.Run(*cfg)
}
