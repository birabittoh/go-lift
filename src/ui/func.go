package ui

import (
	"net/http"
	"time"

	"github.com/birabittoh/go-lift/src/database"
)

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusSeeOther)
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

func formatDay(day string) string {
	// only return the first three letters of the day
	if len(day) < 3 {
		return day
	}
	return day[:3]
}

func isChecked(day string, routine *database.Routine) bool {
	if routine == nil || routine.Days == nil {
		return false
	}
	for _, d := range routine.Days {
		if d.Name == day {
			return true
		}
	}
	return false
}
