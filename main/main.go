package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	urlshort "urlshort/urlpkg"
)

func main() {
	mux := defaultMux()

	filePath := flag.String("file", "", "Path to the YAML file")
	flag.Parse()
	if *filePath == "" {
		log.Fatal("You must specify a file path using the -file flag.")
	}

	// Read the YAML file specified by the file path flag.
	data, err := os.ReadFile(*filePath)

	if err != nil {
		log.Fatalf("Error reading file %s: %v", *filePath, err)
	}

	fileExtension := strings.ToLower(filepath.Ext(*filePath))

	var parser urlshort.DataParser
	switch fileExtension {
	case ".yaml", ".yml":
		parser = urlshort.YamlParser{}
	case ".json":
		parser = urlshort.JsonParser{}
	default:
		log.Fatalf("Unsupported file type: %s", fileExtension)
	}

	handler, err := urlshort.UniversalHandler(parser, data, mux)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
