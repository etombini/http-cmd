package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/etombini/http-cmd/pkg/config"
)

func TestServerConfig(t *testing.T) {
	config_file := os.Getenv("GOPATH") + "/src/github.com/etombini/http-cmd/test-scripts/http-cmd.yaml"

	cfg := config.New(config_file)

	fmt.Printf("TESTING; Config is %+v\n", cfg)

}
