package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type treeRequest struct {
	FavoriteTree string `json:"favoriteTree"`
}

type treeResponse struct {
	Text string `json:"text"`
}

type treeHandler struct {
	tmpl *template.Template
}

// NewTreeHandler creates a handler for the tree webpage
func NewTreeHandler() http.Handler {
	return &treeHandler{
		tmpl: template.Must(template.ParseFiles("index.html")),
	}
}

func (handler *treeHandler) parseRequest(w http.ResponseWriter, r *http.Request) (*treeRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}
	var data treeRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %v", err)
	}
	return &data, nil
}

func (handler *treeHandler) postHandle(w http.ResponseWriter, r *http.Request) {
	var resp treeResponse
	tree, err := handler.parseRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(tree.FavoriteTree) > 0 {
		resp.Text = fmt.Sprintf("It's nice to know that your favorite tree is a %v.", tree.FavoriteTree)
	} else {
		resp.Text = "Please tell me your favorite tree."
	}
	err = handler.tmpl.Execute(w, resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *treeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path != "/" {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}
		handler.postHandle(w, r)
	default:
		http.Error(w, "method not allowed. Methods allowed: POST", http.StatusMethodNotAllowed)
	}
}

func main() {
	handler := NewTreeHandler()
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
