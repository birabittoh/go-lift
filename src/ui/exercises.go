package ui

import (
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
)

func getExercises(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routineItemId, err := getIDFromPath(r)
		if err != nil {
			showError(w, "Routine item ID is required")
			return
		}

		_, err = db.GetRoutineItemByID(routineItemId)
		if err != nil {
			showError(w, "Failed to retrieve routine item: "+err.Error())
			return
		}

		exercises, err := db.GetExercises()
		if err != nil {
			showError(w, "Failed to retrieve exercises: "+err.Error())
			return
		}

		pageData := &PageData{
			Page:      "exercises",
			Exercises: exercises,
			ID:        routineItemId,
		}

		executeTemplateSafe(w, exercisesPath, pageData)
	}
}

func getExercise(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exerciseId := r.PathValue("exerciseId")
		if exerciseId == "" {
			showError(w, "Exercise ID is required")
			return
		}

		routineItemId, err := getIDFromPath(r)
		if err != nil {
			showError(w, "Routine item ID is required")
			return
		}

		_, err = db.GetRoutineItemByID(routineItemId)
		if err != nil {
			showError(w, "Failed to retrieve routine item: "+err.Error())
			return
		}

		exercise, err := db.GetExerciseByID(exerciseId)
		if err != nil {
			showError(w, "Failed to retrieve exercise: "+err.Error())
			return
		}

		pageData := &PageData{
			Page:      "exercises",
			Exercises: []database.Exercise{*exercise},
			ID:        routineItemId,
			Message:   exerciseId,
		}

		executeTemplateSafe(w, exercisePath, pageData)
	}
}
