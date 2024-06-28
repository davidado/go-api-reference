package netjson

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Parse parses the request body into a struct
func Parse(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func Write(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	Write(w, status, map[string]string{"error": err.Error()})
}
