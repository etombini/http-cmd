package main

import (
	"flag"
	"fmt"

	"github.com/etombini/http-cmd/pkg/config"
	"github.com/etombini/http-cmd/pkg/server"
)

func main() {

	configFlag := flag.String("config", config.DefaultConfPath, "Configuration file ["+config.DefaultConfPath+"]")
	flag.Parse()
	fmt.Println("config: ", *configFlag)
	cfg := config.New(*configFlag)

	server.Run(*cfg)
}
