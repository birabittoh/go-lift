package ui

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/birabittoh/go-lift/src/database"
)

func postSetNew(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		itemID, err := getIDFromPath(r)
		if err != nil {
			showError(w, "Invalid exercise item ID: "+err.Error())
			return
		}

		item, err := db.GetExerciseItemByID(itemID)
		if err != nil {
			showError(w, "Exercise item not found")
			return
		}

		_, err = db.NewSet(item)
		if err != nil {
			showError(w, "Failed to create new set: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", item.RoutineItem.RoutineID))
	}
}

func postSet(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		setID, err := getIDFromPath(r)
		if err != nil {
			showError(w, "Invalid set ID: "+err.Error())
			return
		}

		set, err := db.GetSetByID(setID)
		if err != nil {
			showError(w, "Set not found")
			return
		}

		repsStr := r.FormValue("reps")
		if repsStr != "" {
			reps, err := strconv.Atoi(repsStr)
			if err != nil {
				showError(w, "Invalid reps: "+err.Error())
				return
			}

			set.Reps = uint(reps)
		} else {
			set.Reps = 0
		}

		weightStr := r.FormValue("weight")
		if weightStr != "" {
			set.Weight, err = strconv.ParseFloat(weightStr, 64)
			if err != nil {
				showError(w, "Invalid weight: "+err.Error())
				return
			}
		} else {
			set.Weight = 0.0
		}

		durationStr := r.FormValue("duration")
		if durationStr != "" {
			duration, err := strconv.Atoi(durationStr)
			if err != nil {
				showError(w, "Invalid duration: "+err.Error())
				return
			}
			set.Duration = uint(duration)
		} else {
			set.Duration = 0
		}

		err = db.UpdateSet(set)
		if err != nil {
			showError(w, "Failed to update set: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", set.ExerciseItem.RoutineItem.RoutineID))
	}
}

func postSetDelete(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		setID, err := getIDFromPath(r)
		if err != nil {
			showError(w, "Invalid set ID: "+err.Error())
			return
		}

		set, err := db.GetSetByID(setID)
		if err != nil {
			showError(w, "Set not found")
			return
		}

		routineID, err := db.DeleteSet(set)
		if err != nil {
			showError(w, "Failed to delete set: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", routineID))
	}
}
