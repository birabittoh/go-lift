package ui

import (
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
)

func getWorkouts(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		executeTemplateSafe(w, workoutsPath, &PageData{Page: "workouts"})
	}
}
