package ui

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/birabittoh/go-lift/src/database"
	g "github.com/birabittoh/go-lift/src/globals"
)

func postAddExercise(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		routineItemId, err := g.GetIDFromPath(r)
		if err != nil {
			showError(w, "Invalid routine item ID: "+err.Error())
			return
		}

		exerciseId := r.PathValue("exerciseId")
		if exerciseId == "" {
			showError(w, "Exercise ID is required")
			return
		}

		routineItem, err := db.GetRoutineItemByID(routineItemId)
		if err != nil {
			showError(w, "Failed to retrieve routine item: "+err.Error())
			return
		}

		exercise, err := db.GetExerciseByID(exerciseId)
		if err != nil {
			showError(w, "Failed to retrieve exercise: "+err.Error())
			return
		}

		_, err = db.AddExerciseToRoutineItem(routineItem, exercise)
		if err != nil {
			showError(w, "Failed to add exercise to routine item: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", routineItem.RoutineID))
	}
}

func postExerciseItemDelete(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		itemID, err := g.GetIDFromPath(r)
		if err != nil {
			showError(w, "Invalid exercise item ID: "+err.Error())
			return
		}

		item, err := db.GetExerciseItemByID(itemID)
		if err != nil {
			showError(w, "Exercise item not found")
			return
		}

		routineID, err := db.DeleteExerciseItem(item)
		if err != nil {
			showError(w, "Failed to delete exercise item: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", routineID))
	}
}

func postExerciseItem(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		itemID, err := g.GetIDFromPath(r)
		if err != nil {
			showError(w, "Invalid exercise item ID: "+err.Error())
			return
		}

		item, err := db.GetExerciseItemByID(itemID)
		if err != nil {
			showError(w, "Exercise item not found")
			return
		}

		restTimeStr := r.FormValue("restTime")
		if restTimeStr == "" {
			showError(w, "Rest time is required")
			return
		}

		restTime, err := strconv.Atoi(restTimeStr)
		if err != nil {
			showError(w, "Invalid rest time: "+err.Error())
			return
		}

		item.RestTime = uint(restTime)

		err = db.UpdateExerciseItem(item)
		if err != nil {
			showError(w, "Failed to update exercise item: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", item.RoutineItem.RoutineID))
	}
}
