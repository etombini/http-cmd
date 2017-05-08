package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/etombini/http-cmd/pkg/config"
	"github.com/etombini/http-cmd/pkg/server"
)

var v string

// Version returns the version of this application (SemVer format)
// It depends on -ldflags at build time :
// -ldflags "-X github.com/etombini/http-cmd.v=v1.2.3"
// If not set during the build, it defaults to v0.0.0
func version() string {
	if v == "" {
		v = "v0.0.0"
	}
	return v
}

func main() {

	versionFlag := flag.Bool("version", false, "Get version")
	configFlag := flag.String("config", config.DefaultConfPath, "Configuration file ["+config.DefaultConfPath+"]")
	flag.Parse()

	if *versionFlag {
		v := version()
		fmt.Printf("Version %s", v)
		return
	}

	fmt.Println("config: ", *configFlag)
	cfg, err := config.New(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	server.Run(*cfg)
}
