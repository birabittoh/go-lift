package ui

import (
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
)

func getExercises(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exercises, err := db.GetExercises()
		if err != nil {
			showError(w, "Failed to retrieve exercises: "+err.Error())
			return
		}

		pageData := &PageData{
			Page:      "exercises",
			Exercises: exercises,
		}

		executeTemplateSafe(w, exercisesPath, pageData)
	}
}

func getExercise(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			showError(w, "Exercise ID is required")
			return
		}

		exercise, err := db.GetExerciseByID(id)
		if err != nil {
			showError(w, "Failed to retrieve exercise: "+err.Error())
			return
		}

		pageData := &PageData{
			Page:      "exercises",
			Exercises: []database.Exercise{*exercise},
		}

		executeTemplateSafe(w, exercisePath, pageData)
	}
}
