package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"urlshort"
)

func main() {

	var yaml_file string
	flag.StringVar(
		&yaml_file,
		"yaml_file",
		"data/pathsToUrls.yaml",
		"YAML file that maps a path to an HTTP address for redirecting",
	)

	var json_file string
	flag.StringVar(
		&json_file,
		"json_file",
		"data/pathsToUrls.json",
		"JSON file that maps a path to an HTTP address for redirecting",
	)

	flag.Parse()

	// Read in data from JSON file.
	jsn, err := ioutil.ReadFile(json_file)
	if err != nil {
		log.Fatal(err)
	}

	// Read in data from YAML file.
	yml, err := ioutil.ReadFile(yaml_file)
	if err != nil {
		log.Fatal(err)
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(yml, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the YAMLHandler as the
	// fallback
	jsonHandler, err := urlshort.JSONHandler(jsn, yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
