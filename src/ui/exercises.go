package ui

import (
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
	g "github.com/birabittoh/go-lift/src/globals"
)

func getExercises(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageData, err := getPageData(db, "exercises")
		if err != nil {
			showError(w, "Failed to retrieve page data: "+err.Error())
			return
		}

		pageData.ID, err = g.GetIDFromPath(r)
		if err != nil {
			showError(w, "Routine item ID is required")
			return
		}

		_, err = db.GetRoutineItemByID(pageData.ID)
		if err != nil {
			showError(w, "Failed to retrieve routine item: "+err.Error())
			return
		}

		pageData.Exercises, err = db.GetExercises()
		if err != nil {
			showError(w, "Failed to retrieve exercises: "+err.Error())
			return
		}

		executeTemplateSafe(w, exercisesPath, pageData)
	}
}

func getExercise(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageData, err := getPageData(db, "exercises")
		if err != nil {
			showError(w, "Failed to retrieve page data: "+err.Error())
			return
		}

		pageData.Message = r.PathValue("exerciseId")
		if pageData.Message == "" {
			showError(w, "Exercise ID is required")
			return
		}

		pageData.ID, err = g.GetIDFromPath(r)
		if err != nil {
			showError(w, "Routine item ID is required")
			return
		}

		_, err = db.GetRoutineItemByID(pageData.ID)
		if err != nil {
			showError(w, "Failed to retrieve routine item: "+err.Error())
			return
		}

		exercise, err := db.GetExerciseByID(pageData.Message)
		if err != nil {
			showError(w, "Failed to retrieve exercise: "+err.Error())
			return
		}
		pageData.Exercises = []database.Exercise{*exercise}

		executeTemplateSafe(w, exercisePath, pageData)
	}
}
