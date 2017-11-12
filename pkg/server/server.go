package server

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/etombini/http-cmd/pkg/config"
)

func getHandler(config config.Config) http.Handler {
	m := http.NewServeMux()

	ch := catalogHandlerGenerator(config)
	for i := range ch {
		m.HandleFunc(*ch[i].pattern, *ch[i].handler)
	}

	eh := execHandlerGenerator(config)
	for i := range eh {
		m.HandleFunc(*eh[i].pattern, *eh[i].handler)
	}

	return m
}

func getServer(config config.Config) *http.Server {

	var s http.Server
	s.Addr = config.Server.Address
	s.ReadHeaderTimeout = time.Second * 3
	s.WriteTimeout = time.Second * time.Duration(config.Server.Timeout+5)
	s.Handler = getHandler(config)
	return &s
}

// Run starts the server using proper configuration
func Run(config config.Config) {
	server := getServer(config)

	listener, err := net.Listen("tcp", config.Server.Address+":"+strconv.Itoa(int(config.Server.Port)))
	if err != nil {
		os.Exit(1)
	}

	server.Serve(listener)

}
