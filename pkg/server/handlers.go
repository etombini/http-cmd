package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/etombini/http-cmd/pkg/config"
	"github.com/etombini/http-cmd/pkg/hangman"
)

type execHandler struct {
	pattern *string
	handler *func(http.ResponseWriter, *http.Request)
}

type catalogHandler execHandler

// execHandlerGenerator returns a list of struct execHandler.
// Each struct contains an URL and a function which is a http.Handler
func execHandlerGenerator(config config.Config) []execHandler {
	ehs := make([]execHandler, 0)
	for i := range config.Categories {
		for j := range config.Categories[i].Execs {
			pattern := new(string)
			*pattern = config.Server.ExecPrefix + config.Categories[i].Name + "/" + config.Categories[i].Execs[j].Name
			command := new(string)
			*command = config.Categories[i].Execs[j].Command
			timeout := new(uint32)
			*timeout = config.Categories[i].Execs[j].Timeout
			handler := new(func(http.ResponseWriter, *http.Request))

			// Generating the Handler func
			*handler = func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
					return
				}
				if r.URL.Path != *pattern {
					http.NotFound(w, r)
					fmt.Fprintf(os.Stderr, "Invalid URL for command execution (got %s expecting %s)\n", r.URL.Path, *pattern)
					return
				}
				h := hangman.Reaper(*command, *timeout)
				js, err := json.Marshal(h)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					fmt.Fprintf(os.Stderr, "Error while converting execution result to json: %+v", h)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
			eh := execHandler{pattern, handler}
			ehs = append(ehs, eh)
		}
	}
	return ehs
}

type exec4JSON struct {
	Name        string
	Description string
	Command     string
	Timeout     uint32
}

type execCatalog4JSON struct {
	Name        string
	Description string
	Execs       []exec4JSON
}

type catalog4JSON struct {
	Name        string
	Description string
}

// catalogHandlerGenerator returns a list of struct catalogHandler.
// Each struct contains an URL and a function which is a http.Handler
func catalogHandlerGenerator(config config.Config) []catalogHandler {
	chs := make([]catalogHandler, 0)

	c4j := make([]catalog4JSON, 0)
	for i := range config.Categories {
		c := catalog4JSON{config.Categories[i].Name, config.Categories[i].Description}
		c4j = append(c4j, c)
	}
	cPattern := config.Server.CatalogPrefix

	// Generating the Handler func for the first catalog level
	cHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != cPattern {
			http.NotFound(w, r)
			fmt.Fprintf(os.Stderr, "Invalid URL for command execution (%s)\n", r.URL.Path)
			return
		}
		js, err := json.Marshal(c4j)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Fprintf(os.Stderr, "Error while generating global catalog\n")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
	ch := catalogHandler{&cPattern, &cHandler}
	chs = append(chs, ch)

	for i := range config.Categories {
		ecPattern := config.Server.CatalogPrefix + config.Categories[i].Name
		e4j := make([]exec4JSON, 0)
		for j := range config.Categories[i].Execs {
			e := exec4JSON{config.Categories[i].Execs[j].Name,
				config.Categories[i].Execs[j].Description,
				config.Categories[i].Execs[j].Command,
				config.Categories[i].Execs[j].Timeout}
			e4j = append(e4j, e)
		}

		// Generating the Handler func for each catalog category
		ecHandler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Path != ecPattern {
				http.NotFound(w, r)
				fmt.Fprintf(os.Stderr, "Invalid URL for exec catalog (%s)\n", r.URL.Path)
				return
			}
			js, err := json.Marshal(e4j)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Fprintf(os.Stderr, "Error while generating execs catalog\n")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
		ech := catalogHandler{&ecPattern, &ecHandler}
		chs = append(chs, ech)
	}
	return chs
}
