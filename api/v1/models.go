package v1

import (
	"encoding/json"
	"net/http"
)

type ResourceType string

var (
	NoResource   ResourceType
	PostResource ResourceType = "post"
	UserResource ResourceType = "user"
)

type Claim struct {
	Code         string       `json:"code"`
	ResourceType ResourceType `json:"resource_type"`
	Created      int64        `json:"created"`
	Redeemed     int64        `json:"redeemed"`
}

func ServeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
