package functions

import (
	"net/http"
	"strconv"
	"strings"
)

func Capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func GetIDFromPath(r *http.Request) (uint, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	return uint(id), err
}
