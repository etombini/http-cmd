package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"

	"errors"

	yaml "gopkg.in/yaml.v2"
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
		Port          uint32 `yaml:"port"`
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

func checkServerDefault(c *Config) error {
	if c.Server.Address == "" {
		fmt.Fprintf(os.Stderr, "Address is not set, defaulting to %s\n", DefaultAddress)
		c.Server.Address = DefaultAddress
	}
	if net.ParseIP(c.Server.Address) == nil {
		fmt.Fprintf(os.Stderr, "Address %s is not a valid IP (v4 or v6) address\n", c.Server.Address)
		return errors.New("Address " + c.Server.Address + " is not a valid IP (v4 or v6) address")
	}
	if c.Server.Port == 0 {
		fmt.Fprintf(os.Stderr, "Port is not set, defaulting to %d\n", DefaultPort)
		c.Server.Port = DefaultPort
	}
	if c.Server.Port < 0 || c.Server.Port > 65535 {
		fmt.Fprintf(os.Stderr, "Port number must be in 1-65535\n")
		return errors.New("Port number must be in 1-65535")
	}
	if c.Server.Timeout <= 0 {
		fmt.Fprintf(os.Stderr, "Timeout is not set, defaulting to %d\n", DefaultTimeout)
		c.Server.Port = DefaultTimeout
	}
	if c.Server.CatalogPrefix == "" {
		fmt.Fprintf(os.Stderr, "Catalog prefix is not set, defaulting to %s\n", DefaultCatalogPrefix)
		c.Server.CatalogPrefix = DefaultCatalogPrefix
	}
	if !strings.HasPrefix(c.Server.CatalogPrefix, "/") {
		c.Server.CatalogPrefix = "/" + c.Server.CatalogPrefix
	}
	if !strings.HasSuffix(c.Server.CatalogPrefix, "/") {
		c.Server.CatalogPrefix = c.Server.CatalogPrefix + "/"
	}
	if c.Server.ExecPrefix == "" {
		fmt.Fprintf(os.Stderr, "Exec prefix is not set, defaulting to %s\n", DefaultExecPrefix)
		c.Server.ExecPrefix = DefaultExecPrefix
	}
	if !strings.HasPrefix(c.Server.ExecPrefix, "/") {
		c.Server.ExecPrefix = "/" + c.Server.ExecPrefix
	}
	if !strings.HasSuffix(c.Server.ExecPrefix, "/") {
		c.Server.ExecPrefix = c.Server.ExecPrefix + "/"
	}
	if c.Server.CatalogPrefix == c.Server.ExecPrefix {
		fmt.Fprintf(os.Stderr, "Exec prefix (%s) and Catalog prefix (%s) can not have the same value\n",
			c.Server.CatalogPrefix,
			c.Server.ExecPrefix)
	}
	return nil
}

func checkCategoryDuplicates(c *Config) error {
	m := make(map[string]bool)
	for i := range c.Categories {
		category := c.Categories[i].Name
		if m[category] {
			fmt.Fprintf(os.Stderr, "Category duplicate found: %s - exiting\n", category)
			return errors.New("Category duplicate found: " + category + " - exiting")
		}
		m[category] = true
	}
	return nil
}

func checkCategoryNames(c *Config) error {
	for i := range c.Categories {
		if strings.Contains(c.Categories[i].Name, "/") {
			fmt.Fprintf(os.Stderr, "Category name (%s) must not contain a \"/\"\n", c.Categories[i].Name)
			return errors.New("Category name (" + c.Categories[i].Name + ") must not contain a \"/\"")
		}
		if c.Categories[i].Name == "" {
			fmt.Fprintf(os.Stderr, "Category name can not be an empty string\n")
			return errors.New("Category name can not be an empty string")
		}
	}
	return nil
}

func normalizeExecsPath(c *Config) error {
	dir, _ := filepath.Split(c.FilePath)
	for i := range c.Categories {
		path := c.Categories[i].ExecsFilePath
		if !strings.HasPrefix(path, "/") {
			path = dir + path
			c.Categories[i].ExecsFilePath = path
		}
	}
	return nil
}

func loadExecs(c *Config) error {
	for i := range c.Categories {
		ePath := c.Categories[i].ExecsFilePath
		if _, err := os.Stat(ePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Exec path %s in category %s does not exist\n",
				ePath, c.Categories[i].Name)
			return errors.New("Exec path " + ePath + " in category " + c.Categories[i].Name + " does not exist")
		}
		config, err := ioutil.ReadFile(ePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not find or open configuration file %s - %s\n", ePath, err.Error())
			return errors.New("Can not find or open configuration file " + ePath + ": " + err.Error())
		}
		eConfig := Execs{}
		if err := yaml.Unmarshal(config, &eConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Can not parse configuration file %s - %s\n", ePath, err.Error())
			return errors.New("Can not parse configuration file " + ePath + ": " + err.Error())
		}
		// check for duplicates
		m := make(map[string]bool)
		for j := range eConfig.Execs {
			name := eConfig.Execs[j].Name
			if m[name] {
				fmt.Fprintf(os.Stderr, "Exec duplicate found (%s) in category %s (%s)\n",
					name, c.Categories[i].Name, c.Categories[i].ExecsFilePath)
				return errors.New("Exec duplicate found (" + name + ") in category " + c.Categories[i].Name + " (" + c.Categories[i].ExecsFilePath + ")")
			}
			m[name] = true
		}
		c.Categories[i].Execs = eConfig.Execs
	}
	return nil
}

func checkExecNames(c *Config) error {
	for i := range c.Categories {
		for j := range c.Categories[i].Execs {
			if strings.Contains(c.Categories[i].Execs[j].Name, "/") {
				fmt.Fprintf(os.Stderr, "Exec name (%s) must not contain a \"/\"\n", c.Categories[i].Execs[j].Name)
				return errors.New("Exec name (" + c.Categories[i].Execs[j].Name + ") must not contain a \"/\"")
			}
			if c.Categories[i].Execs[j].Name == "" {
				fmt.Fprintf(os.Stderr, "Category name can not be an empty string\n")
				return errors.New("Category name can not be an empty string")
			}
		}
	}
	return nil
}

// New return a new Config structure, loaded according to a configuration file
func New(filename string) (*Config, error) {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not find or open configuration file %s\n", filename)
		return nil, err
	}

	cfg := Config{}
	if err := yaml.Unmarshal(config, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Can not parse configuration file %s - %s\n", filename, err.Error())
		return nil, err
	}
	cfg.FilePath = filename

	if err := checkServerDefault(&cfg); err != nil {
		return nil, err
	}
	if err := checkCategoryNames(&cfg); err != nil {
		return nil, err
	}
	if err := checkCategoryDuplicates(&cfg); err != nil {
		return nil, err
	}
	if err := normalizeExecsPath(&cfg); err != nil {
		return nil, err
	}
	if err := loadExecs(&cfg); err != nil {
		return nil, err
	}
	if err := checkExecNames(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
