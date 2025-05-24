package api

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func getIDFromPath(r *http.Request) (uint, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	return uint(id), err
}
