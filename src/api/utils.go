package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, messages ...string) {
	jsonResponse(w, status, map[string]string{"error": strings.Join(messages, " ")})
}
