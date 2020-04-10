package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/pbabbicola/tree-server-test/tree"
)

const (
	Addr      = ":8000"
	IndexFile = "html/index.html"
)

func run() error {
	index, err := template.ParseFiles(IndexFile)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", tree.NewHandler(index))

	return http.ListenAndServe(Addr, mux)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
