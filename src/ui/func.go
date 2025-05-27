package ui

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getIDFromPath(r *http.Request) (uint, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	return uint(id), err
}

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func coalesce(s1 *string, s string) string {
	if s1 != nil && *s1 != "" {
		return *s1
	}
	return s
}

func formatBirthDate(birthdate *time.Time) string {
	if birthdate == nil {
		return ""
	}
	return birthdate.Format("2006-01-02")
}
