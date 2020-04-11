package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/pbabbicola/tree-server-test/tree"
)

const (
	Addr      = "localhost:8000"
	IndexFile = "html/index.html"
)

func run() error {
	index, err := template.ParseFiles(IndexFile)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	log.Printf("listening at %v\n", Addr)
	return http.ListenAndServe(Addr, tree.NewHandler(index))
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
