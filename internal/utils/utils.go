package utils

import (
	"encoding/json"
	"net/http"
)

func ParseRequestBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func CreateResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{}
	response["data"] = data
	json.NewEncoder(w).Encode(response)
}
