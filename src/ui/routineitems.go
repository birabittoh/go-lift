package ui

import (
	"fmt"
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
)

func postRoutineItemNew(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, err := getIDFromPath(r)
		if err != nil {
			showError(w, "Invalid routine ID: "+err.Error())
			return
		}

		routine, err := db.GetRoutineByID(id)
		if err != nil {
			showError(w, "Failed to retrieve routine: "+err.Error())
			return
		}

		item, err := db.NewRoutineItem(routine)
		if err != nil {
			showError(w, "Failed to create new routine item: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/exercises/%d", item.ID))
	}
}

func postRoutineItemDelete(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		itemID, err := getIDFromPath(r)
		if err != nil {
			showError(w, "Invalid routine item ID: "+err.Error())
			return
		}

		item, err := db.GetRoutineItemByID(itemID)
		if err != nil {
			showError(w, "Routine item not found")
			return
		}

		routineID, err := db.DeleteRoutineItem(item)
		if err != nil {
			showError(w, "Failed to delete routine item: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", routineID))
	}
}
