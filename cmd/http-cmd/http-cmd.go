package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/etombini/http-cmd/pkg/config"
	"github.com/etombini/http-cmd/pkg/server"
)

var (
	v string
	b string
)

func version() string {
	if v == "" {
		v = "v0.0.0"
	}
	return v
}

func build() string {
	if b == "" {
		b = "Unknown"
	}
	return b
}

func main() {

	versionFlag := flag.Bool("version", false, "Get version")
	configFlag := flag.String("config", config.DefaultConfPath, "Configuration file ["+config.DefaultConfPath+"]")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version %s\n", version())
		fmt.Printf("Build: %s\n", build())
		return
	}

	fmt.Println("config: ", *configFlag)
	cfg, err := config.New(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	server.Run(*cfg)
}
