package v1

import (
	"encoding/json"
	"net/http"
)

type Grant struct {
	Code         string `json:"code"`
	ResourceType string `json:"resource_type"`
	Created      int64  `json:"created"`
}

func ServeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
