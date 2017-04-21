package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	// DefaultConfPath is the default configuration file path
	DefaultConfPath string = "/etc/http-cmd/http-cmd.conf"
	// DefaultCatalogPath is the default path for catalog/categories
	DefaultCatalogPath string = "/etc/http-cmd/http-cmd-catalog.conf"
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
	// DefaultRunPrefix is the default URL prefix to reach command execution
	DefaultRunPrefix string = "/run/"
	// LoggerName is the default logger name for this package
	LoggerName string = "config"
)

// Server handles server configuration and the catalog of commands
// to be executed
type Server struct {
	FilePath      string
	Port          int
	Address       string
	Timeout       int
	CatalogPrefix string
	RunPrefix     string
	Categories    []Category
}

// Category handles a category configuration and affiliated commands
type Category struct {
	ExecFilePath string
	Name         string
	Description  string
	Execs        []Exec
}

// Exec handles a command to be executed by the server
type Exec struct {
	Name        string
	Command     string
	Timeout     int
	Description string
}

//GetConfig returns the server configuration, including categories and executer
func GetConfig(cfgfile string) *Server {
	if cfgfile == "" {
		fmt.Printf("Empty configuration file path, using default %s\n", DefaultConfPath)
		cfgfile = DefaultConfPath
	}
	if _, err := os.Stat(cfgfile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Configuration file %s not found\n", cfgfile)
		os.Exit(1)
	}
	dirConfig, _ := filepath.Split(cfgfile)
	var serverConfig Server
	serverConfig.FilePath = cfgfile
	serverConfig.Port = DefaultPort
	serverConfig.Address = DefaultAddress
	serverConfig.Timeout = DefaultTimeout
	serverConfig.CatalogPrefix = DefaultCatalogPrefix
	serverConfig.RunPrefix = DefaultRunPrefix
	serverConfig.Categories = make([]Category, 0)

	serverConfig.FilePath = cfgfile

	vServerConfig := viper.New()
	vServerConfig.SetConfigFile(serverConfig.FilePath)
	vServerConfig.SetConfigType("yaml")
	if err := vServerConfig.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	serverSubConfig := vServerConfig.Sub("server")

	if v := serverSubConfig.GetInt("port"); v != 0 {
		serverConfig.Port = v
	}
	if v := serverSubConfig.GetString("address"); v != "" {
		serverConfig.Address = v
	}
	if v := serverSubConfig.GetInt("timeout"); v != 0 {
		serverConfig.Timeout = v
	}
	if v := serverSubConfig.GetString("Server.CatalogPrefix"); v != "" {
		if !strings.HasPrefix(v, "/") {
			v = "/" + v
		}
		if !strings.HasSuffix(v, "/") {
			v = v + "/"
		}
		serverConfig.CatalogPrefix = v
	}
	if v := serverSubConfig.GetString("Server.RunPrefix"); v != "" {
		if !strings.HasPrefix(v, "/") {
			v = "/" + v
		}
		if !strings.HasSuffix(v, "/") {
			v = v + "/"
		}
		if serverConfig.CatalogPrefix == v {
			fmt.Println("Catalog prefix can not be the same as Run prefix")
			os.Exit(1)
		}
		serverConfig.RunPrefix = v
	}

	categoriesCfg := vServerConfig.Get("categories").([]interface{})
	for i := range categoriesCfg {
		var c Category
		c.Description = DefaultDescription
		c.Execs = make([]Exec, 0)
		for k, v := range categoriesCfg[i].(map[interface{}]interface{}) {
			switch k.(string) {
			case "name":
				c.Name = v.(string)
			case "description":
				c.Description = v.(string)
			case "execs":
				if strings.HasSuffix(v.(string), "/") {
					c.ExecFilePath = v.(string)
				} else {
					c.ExecFilePath = dirConfig + v.(string)
				}
			default:
				fmt.Fprintf(os.Stderr, "Unknown category configuration key %s (%s)\n", k.(string), v.(string))
			}
		}
		if c.Name == "" {
			fmt.Fprint(os.Stderr, "A category must have a name\n")
			os.Exit(1)
		}
		if c.ExecFilePath == "" {
			fmt.Fprintf(os.Stderr, "A category (%s) must refer to a file containing execs\n", c.Name)
			os.Exit(1)
		}
		serverConfig.Categories = append(serverConfig.Categories, c)
	}

	for i := range serverConfig.Categories {
		c := &serverConfig.Categories[i]
		if _, err := os.Stat(c.ExecFilePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Exec configuration file %s not found\n", c.ExecFilePath)
			os.Exit(1)
		}
		vExecConfig := viper.New()
		vExecConfig.SetConfigFile(c.ExecFilePath)
		vExecConfig.SetConfigType("yaml")
		if err := vExecConfig.ReadInConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		execsCfg := vExecConfig.Get("execs").([]interface{})
		for i := range execsCfg {
			var e Exec
			e.Description = DefaultDescription
			e.Timeout = DefaultTimeout

			for k, v := range execsCfg[i].(map[interface{}]interface{}) {
				switch k.(string) {
				case "name":
					e.Name = v.(string)
				case "description":
					e.Description = v.(string)
				case "timeout":
					e.Timeout = v.(int)
				case "command":
					e.Command = v.(string)
				default:
					fmt.Fprintf(os.Stderr, "Unknown exec configuration key %s (%s)\n", k.(string), v.(string))
				}
			}
			if e.Name == "" {
				fmt.Fprint(os.Stderr, "An exec must have a name\n")
				os.Exit(1)
			}
			if e.Command == "" {
				fmt.Fprintf(os.Stderr, "An exec (%s) must have a command to execute\n", e.Name)
				os.Exit(1)
			}
			c.Execs = append(c.Execs, e)
		}
	}

	return &serverConfig
}
