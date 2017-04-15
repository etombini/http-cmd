package config_test

import (
	"fmt"
	"github.com/etombini/http-cmd/pkg/config"
	"os"
	"testing"
)

func TestServerConfig(t *testing.T) {
	config_file := os.Getenv("GOPATH") + "/src/github.com/etombini/http-cmd/test-scripts/http-cmd.conf"

	cfg := config.GetConfig(config_file)

	fmt.Printf("TESTING; Config is %+v\n", cfg)

}
