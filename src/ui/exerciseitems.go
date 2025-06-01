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

func postExerciseItemsDelete(db *database.Database) http.HandlerFunc {
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

func postExerciseItemsUp(db *database.Database) http.HandlerFunc {
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

		if item.OrderIndex == 0 {
			showError(w, "Cannot move the first exercise item up")
			return
		}

		err = db.MoveExerciseItem(item, true) // true for moving up
		if err != nil {
			showError(w, "Failed to move exercise item up: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", item.RoutineItem.RoutineID))
	}
}

func postExerciseItemsDown(db *database.Database) http.HandlerFunc {
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

		if item.OrderIndex == len(item.RoutineItem.ExerciseItems)-1 {
			showError(w, "Cannot move the last exercise item down")
			return
		}

		err = db.MoveExerciseItem(item, false) // false for moving down
		if err != nil {
			showError(w, "Failed to move exercise item down: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", item.RoutineItem.RoutineID))
	}
}

func postExerciseItems(db *database.Database) http.HandlerFunc {
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

		if err := r.ParseForm(); err != nil {
			showError(w, "Failed to parse form: "+err.Error())
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
		item.Notes = r.FormValue("notes")

		err = db.UpdateExerciseItem(item)
		if err != nil {
			showError(w, "Failed to update exercise item: "+err.Error())
			return
		}

		for i, set := range item.Sets {
			prefix := fmt.Sprintf("sets[%d]", i)

			repsStr := r.FormValue(prefix + "[reps]")
			if repsStr != "" {
				reps, err := strconv.Atoi(repsStr)
				if err != nil {
					showError(w, fmt.Sprintf("Invalid reps for set %d: %v", i+1, err))
					return
				}
				uReps := uint(reps)
				set.Reps = &uReps
			} else {
				set.Reps = nil
			}

			weightStr := r.FormValue(prefix + "[weight]")
			if weightStr != "" {
				weight, err := strconv.ParseFloat(weightStr, 64)
				if err != nil {
					showError(w, fmt.Sprintf("Invalid weight for set %d: %v", i+1, err))
					return
				}
				set.Weight = &weight
			} else {
				set.Weight = nil
			}

			durationStr := r.FormValue(prefix + "[duration]")
			if durationStr != "" {
				duration, err := strconv.Atoi(durationStr)
				if err != nil {
					showError(w, fmt.Sprintf("Invalid duration for set %d: %v", i+1, err))
					return
				}
				uDuration := uint(duration)
				set.Duration = &uDuration
			} else {
				set.Duration = nil
			}

			if err := db.UpdateSet(&set); err != nil {
				showError(w, fmt.Sprintf("Failed to update set %d: %v", i+1, err))
				return
			}
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", item.RoutineItem.RoutineID))
	}
}
