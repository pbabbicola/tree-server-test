package tree

import (
	"fmt"
	"net/http"
	"text/template"
)

type treeResponse struct {
	Text string `json:"text"`
}

type Handler struct {
	index *template.Template
}

// NewHandler creates a handler for the tree webpage
func NewHandler(index *template.Template) *Handler {
	return &Handler{
		index: index,
	}
}

func getFavoriteTree(r *http.Request) string {
	t, ok := r.URL.Query()["favoriteTree"]
	if !ok || len(t) == 0 {
		return ""
	}

	return t[0] // discard any other listed trees, we only need one
}

func formatResponse(tree string) string {
	if tree != "" {
		return fmt.Sprintf("It's nice to know that your favorite tree is a %v.", tree)
	}

	return "Please tell me your favorite tree."
}

func (handler *Handler) renderResponse(w http.ResponseWriter, text string) {
	resp := treeResponse{
		Text: text,
	}

	err := handler.index.Execute(w, resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *Handler) getHandle(w http.ResponseWriter, r *http.Request) {
	tree := getFavoriteTree(r)
	text := formatResponse(tree)
	handler.renderResponse(w, text)
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/" {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}

		handler.getHandle(w, r)

	default:
		http.Error(w, "method not allowed. Methods allowed: GET", http.StatusMethodNotAllowed)
	}
}
