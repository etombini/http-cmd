package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultConfPath is the default configuration file path
	DefaultConfPath string = "/etc/http-cmd/http-cmd.conf"
	// DefaultPort is the default port for the server to listen
	DefaultPort int = 5050
	// DefaultAddress is the default address the server is binding
	DefaultAddress string = "127.0.0.1"
	// DefaultTimeout is the default timeout for command execution
	DefaultTimeout int = 5
	// DefaultDescription is the default description for categories and command
	DefaultDescription string = "No description provided"
	// DefaultCatalogPrefix is the default URL prefix to reach and show the catalog
	DefaultCatalogPrefix string = "/catalog/"
	// DefaultExecPrefix is the default URL prefix to reach command execution
	DefaultExecPrefix string = "/run/"
	// LoggerName is the default logger name for this package
	LoggerName string = "config"
)

// Config is a structure representing the global application configuration
type Config struct {
	Server struct {
		Address       string `yaml:"address"`
		Port          int    `yaml:"port"`
		Timeout       int    `yaml:"timeout"`
		CatalogPrefix string `yaml:"catalog_prefix"`
		ExecPrefix    string `yaml:"exec_prefix"`
	}

	FilePath   string
	Categories []Category `yaml:"categories"`
}

// Category is a structure handling category configuration
type Category struct {
	Name          string `yaml:"name"`
	Description   string `yaml:"description"`
	ExecsFilePath string `yaml:"path"`
	Execs         []Exec
}

// Execs is a list of Exec used only to comply to yaml unmarshalling
type Execs struct {
	Execs []Exec `yaml:"execs"`
}

// Exec is a structure handling exec configuration
type Exec struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Command     string `yaml:"command"`
	Timeout     int    `yaml:"timeout"`
}

func (c *Config) checkServerDefault() {
	if c.Server.Address == "" {
		fmt.Fprintf(os.Stderr, "Address is not set, defaulting to %s\n", DefaultAddress)
		c.Server.Address = DefaultAddress
	}
	if net.ParseIP(c.Server.Address) == nil {
		fmt.Fprintf(os.Stderr, "Address %s is not a valid IP (v4 or v6) address\n", c.Server.Address)
		os.Exit(1)
	}
	if c.Server.Port <= 0 {
		fmt.Fprintf(os.Stderr, "Port is not set, defaulting to %d\n", DefaultPort)
		c.Server.Port = DefaultPort
	}
	if c.Server.Timeout <= 0 {
		fmt.Fprintf(os.Stderr, "Timeout is not set, defaulting to %d\n", DefaultTimeout)
		c.Server.Port = DefaultTimeout
	}
	if c.Server.CatalogPrefix == "" {
		fmt.Fprintf(os.Stderr, "Catalog prefix is not set, defaulting to %s\n", DefaultCatalogPrefix)
		c.Server.CatalogPrefix = DefaultCatalogPrefix
	}
	if c.Server.ExecPrefix == "" {
		fmt.Fprintf(os.Stderr, "Exec prefix is not set, defaulting to %s\n", DefaultExecPrefix)
		c.Server.ExecPrefix = DefaultExecPrefix
	}
	if c.Server.CatalogPrefix == c.Server.ExecPrefix {
		fmt.Fprintf(os.Stderr, "Exec prefix (%s) and Catalog prefix (%s) can not have the same value\n",
			c.Server.CatalogPrefix,
			c.Server.ExecPrefix)
	}
}

func (c *Config) checkCategoryDuplicates() {
	m := make(map[string]bool)
	for i := range c.Categories {
		category := c.Categories[i].Name
		if m[category] {
			fmt.Fprintf(os.Stderr, "Category duplicate found: %s - exiting\n", category)
			os.Exit(1)
		}
		m[category] = true
	}
}

func (c *Config) normalizeExecsPath() {
	dir, _ := filepath.Split(c.FilePath)
	for i := range c.Categories {
		path := c.Categories[i].ExecsFilePath
		if !strings.HasPrefix(path, "/") {
			path = dir + path
			c.Categories[i].ExecsFilePath = path
		}
	}
}

func (c *Config) loadExecs() {
	for i := range c.Categories {
		ePath := c.Categories[i].ExecsFilePath
		if _, err := os.Stat(ePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Exec path %s in category %s does not exist",
				ePath, c.Categories[i].Name)
			os.Exit(1)
		}
		config, err := ioutil.ReadFile(ePath)
		if err != nil {
			panic(err)
		}
		eConfig := Execs{}
		if err := yaml.Unmarshal(config, &eConfig); err != nil {
			panic(err)
		}
		//TODO: check for empty name -> Create 2 functions : checkExecsDuplicates & validateExecFormat (and Category as well)
		// check for duplicates
		m := make(map[string]bool)
		for j := range eConfig.Execs {
			name := eConfig.Execs[j].Name
			if m[name] {
				fmt.Fprintf(os.Stderr, "Exec duplicate found (%s) in category %s (%s)\n",
					name, c.Categories[i].Name, c.Categories[i].ExecsFilePath)
				os.Exit(1)
			}
			m[name] = true
		}
		c.Categories[i].Execs = eConfig.Execs
	}
}

// New return a new Config structure, loaded according to a configuration file
func New(filename string) *Config {
	// TODO: change this ckeck to a better one (using the default configuration path)
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	cfg := Config{}
	if err := yaml.Unmarshal(config, &cfg); err != nil {
		panic(err)
	}
	cfg.FilePath = filename

	cfg.checkServerDefault()
	cfg.checkCategoryDuplicates()
	cfg.normalizeExecsPath()

	cfg.loadExecs()

	return &cfg
}
