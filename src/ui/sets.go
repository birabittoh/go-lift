package ui

import (
	"fmt"
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
	g "github.com/birabittoh/go-lift/src/globals"
)

func postSetNew(db *database.Database) http.HandlerFunc {
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

		_, err = db.NewSet(item)
		if err != nil {
			showError(w, "Failed to create new set: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", item.RoutineItem.RoutineID))
	}
}

func postSetDelete(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		setID, err := g.GetIDFromPath(r)
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
