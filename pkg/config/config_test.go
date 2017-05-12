package config_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/etombini/http-cmd/pkg/config"
)

func TestServerConfig(t *testing.T) {
	config_file := os.Getenv("GOPATH") + "/src/github.com/etombini/http-cmd/test-scripts/http-cmd.yaml"

	cfg, _ := config.New(config_file)

	fmt.Printf("TESTING; Config is %+v\n", cfg)

}

func TestConfigServerDefault(t *testing.T) {
	configFile := os.Getenv("GOPATH") + "/src/github.com/etombini/http-cmd/test-scripts/config/missing-default/http-cmd.yaml"

	cfg, err := config.New(configFile)
	if err != nil {
		t.Error("TestConfigServerDefault: Error while creating Config: " + err.Error())
		return
	}

	if cfg.Server.Address != config.DefaultAddress {
		t.Error("TestConfigServerDefault: Default server address is not "+config.DefaultAddress+" :", cfg.Server.Address)
	}
	if cfg.Server.Port != config.DefaultPort {
		t.Error("TestConfigServerDefault: Default server port is not "+strconv.Itoa(int(config.DefaultPort))+": ", cfg.Server.Port)
	}
	if cfg.Server.Timeout != config.DefaultTimeout {
		t.Error("TestConfigServerDefault: Default server timeout is not "+strconv.Itoa(int(config.DefaultTimeout))+": ", cfg.Server.Timeout)
	}
	if cfg.Server.CatalogPrefix != config.DefaultCatalogPrefix {
		t.Error("TestConfigServerDefault: Default server catalog prefix is not "+config.DefaultCatalogPrefix+": ", cfg.Server.CatalogPrefix)
	}
	if cfg.Server.ExecPrefix != config.DefaultExecPrefix {
		t.Error("TestConfigServerDefault: Default server catalog prefix is not "+config.DefaultExecPrefix+": ", cfg.Server.CatalogPrefix)
	}

}

func TestConfigServerAddressOutOfBound(t *testing.T) {
	configFile := os.Getenv("GOPATH") + "/src/github.com/etombini/http-cmd/test-scripts/config/out-of-bound/http-cmd-address-01.yaml"
	cfg, err := config.New(configFile)
	if err == nil {
		t.Error("TestConfigServerDefault: Missing error for bad address " + cfg.Server.Address)
		return
	}

}

func TestConfigServerPortOutOfBound(t *testing.T) {
	{
		configFile := os.Getenv("GOPATH") + "/src/github.com/etombini/http-cmd/test-scripts/config/out-of-bound/http-cmd-port-01.yaml"
		cfg, err := config.New(configFile)
		if err == nil {
			t.Error("TestConfigServerDefault: Missing error for bad address " + cfg.Server.Address)
			return
		}
	}
	{
		configFile := os.Getenv("GOPATH") + "/src/github.com/etombini/http-cmd/test-scripts/config/out-of-bound/http-cmd-port-02.yaml"
		cfg, err := config.New(configFile)
		if err == nil {
			t.Error("TestConfigServerDefault: Missing error for bad address " + cfg.Server.Address)
			return
		}
	}
	{
		configFile := os.Getenv("GOPATH") + "/src/github.com/etombini/http-cmd/test-scripts/config/out-of-bound/http-cmd-port-03.yaml"
		cfg, err := config.New(configFile)
		if err == nil {
			t.Error("TestConfigServerDefault: Missing error for bad address " + cfg.Server.Address)
			return
		}
	}
}
