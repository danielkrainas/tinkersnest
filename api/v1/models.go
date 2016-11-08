package v1

import (
	"encoding/json"
	"net/http"
)

type CreatePostRequest struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Body  string `json:"body"`
}

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Body    string `json:"body"`
	Created int64  `json:"created"`
}

func ServeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
