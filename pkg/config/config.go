package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultConfPath      string = "/etc/http-cmd/http-cmd.conf"
	DefaultCatalogPath   string = "/etc/http-cmd/http-cmd-catalog.conf"
	DefaultPort          int    = 5050
	DefaultAddress       string = "127.0.0.1"
	DefaultTimeout       int    = 5
	DefaultDescription   string = "No description provided"
	DefaultPrefixCatalog string = "/catalog/"
	DefaultPrefixRun     string = "/run/"
	LoggerName           string = "config"
)

type ServerConfig struct {
	FilePath      string
	CategoryPath  string
	Port          int
	Address       string
	Timeout       int
	PrefixCatalog string
	PrefixRun     string
	Exec          map[string]CategoryConfig
}

type CategoryConfig struct {
	Name        string
	ExecsPath   string
	Description string
	Execs       map[string]ExecConfig
}

type ExecConfig struct {
	Name        string
	Command     string
	Description string
	Timeout     int
}

func GetConfig(cfgfile string) *ServerConfig {
	if cfgfile == "" {
		fmt.Printf("Empty configuration file path, using default %s\n", DefaultConfPath)
		cfgfile = DefaultConfPath
	}
	if _, err := os.Stat(cfgfile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Configuration file %s not found\n", cfgfile)
		os.Exit(1)
	}
	var sc ServerConfig
	sc.FilePath = cfgfile
	sc.Port = DefaultPort
	sc.Address = DefaultAddress
	sc.Timeout = DefaultTimeout
	sc.PrefixCatalog = DefaultPrefixCatalog
	sc.PrefixRun = DefaultPrefixRun
	if sc.Exec == nil {
		sc.Exec = make(map[string]CategoryConfig)
	}

	sc.FilePath = cfgfile

	server_cfg := viper.New()
	server_cfg.SetConfigFile(sc.FilePath)
	server_cfg.SetConfigType("toml")
	if err := server_cfg.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if v := server_cfg.GetString("Server.Catalog"); v != "" {
		if strings.HasPrefix(v, "/") {
			sc.CategoryPath = v
		} else {
			dir, _ := filepath.Split(sc.FilePath)
			sc.CategoryPath = dir + v
		}
	}
	if v := server_cfg.GetInt("Server.Port"); v != 0 {
		sc.Port = server_cfg.GetInt("Server.Port")
	}
	if v := server_cfg.GetString("Server.Address"); v != "" {
		sc.Address = server_cfg.GetString("Server.Address")
	}
	if v := server_cfg.GetInt("Server.Timeout"); v != 0 {
		sc.Timeout = server_cfg.GetInt("Server.Timeout")
	}
	if v := server_cfg.GetString("Server.CatalogPrefix"); v != "" {
		if !strings.HasPrefix(v, "/") {
			v = "/" + v
		}
		if !strings.HasSuffix(v, "/") {
			v = v + "/"
		}
		sc.PrefixCatalog = v
	}
	if v := server_cfg.GetString("Server.RunPrefix"); v != "" {
		if !strings.HasPrefix(v, "/") {
			v = "/" + v
		}
		if !strings.HasSuffix(v, "/") {
			v = v + "/"
		}
		if sc.PrefixCatalog == v {
			fmt.Println("Catalog prefix can not be the same as Run prefix")
			os.Exit(1)
		}
		sc.PrefixRun = v
	}

	if _, err := os.Stat(sc.CategoryPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Category configuration file %s not found\n", sc.CategoryPath)
		os.Exit(1)
	}

	category_cfg := viper.New()
	category_cfg.SetConfigFile(sc.CategoryPath)
	category_cfg.SetConfigType("toml")
	if err := category_cfg.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	categories := make(map[string]bool)
	valued_keys := category_cfg.AllKeys()
	for _, key := range valued_keys {
		if strings.HasSuffix(key, ".execs") {
			categories[key[0:len(key)-len(".execs")]] = true
		}
	}

	for category := range categories {
		var cc CategoryConfig
		cc.Name = category
		cc.ExecsPath = ""
		cc.Description = DefaultDescription
		cc.Execs = make(map[string]ExecConfig)

		if v := category_cfg.GetString(category + ".Description"); v != "" {
			cc.Description = v
		}
		if v := category_cfg.GetString(category + ".Execs"); v != "" {
			if strings.HasPrefix(v, "/") {
				cc.ExecsPath = v
			} else {
				dir, _ := filepath.Split(sc.CategoryPath)
				cc.ExecsPath = dir + v
			}
		} else {
			fmt.Printf("Execs file path for category %s is not set\n", category)
			os.Exit(1)
		}
		sc.Exec[category] = cc
	}

	for category := range sc.Exec {
		if _, err := os.Stat(sc.Exec[category].ExecsPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Configuration file %s not found\n", sc.Exec[category].ExecsPath)
			os.Exit(1)
		}

		exec_cfg := viper.New()
		exec_cfg.SetConfigFile(sc.Exec[category].ExecsPath)
		exec_cfg.SetConfigType("toml")
		if err := exec_cfg.ReadInConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		execs := make(map[string]bool)
		valued_keys := exec_cfg.AllKeys()
		for _, key := range valued_keys {
			if strings.HasSuffix(key, ".command") {
				execs[key[0:len(key)-len(".command")]] = true
			}
		}

		for exec := range execs {
			var ec ExecConfig
			ec.Name = exec
			ec.Command = ""
			ec.Description = DefaultDescription
			ec.Timeout = DefaultTimeout

			if v := exec_cfg.GetString(exec + ".Command"); v != "" {
				ec.Command = v
			} else {
				fmt.Fprintln(os.Stderr, "Command for %s not found in %s", exec, sc.Exec[category].ExecsPath)
				os.Exit(1)
			}
			if v := exec_cfg.GetString(exec + ".Description"); v != "" {
				ec.Description = v
			}
			if v := exec_cfg.GetInt(exec + ".Timeout"); v != 0 {
				ec.Timeout = v
			}
			sc.Exec[category].Execs[ec.Name] = ec
		}

	}

	// for category := range sc.Exec {
	// 	for exec := range sc.Exec[category].Execs {
	// 		fmt.Printf("/%s/%s for command %s\n", category, exec, sc.Exec[category].Execs[exec].Command)
	// 	}
	// }

	return &sc
}
