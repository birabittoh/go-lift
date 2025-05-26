package ui

import (
	"strings"
	"time"
)

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
