package ui

import (
	"fmt"
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
	g "github.com/birabittoh/go-lift/src/globals"
)

func postAddRoutineItems(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, err := g.GetIDFromPath(r)
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

func postRoutineItemsUp(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		itemID, err := g.GetIDFromPath(r)
		if err != nil {
			showError(w, "Invalid routine item ID: "+err.Error())
			return
		}

		item, err := db.GetRoutineItemByID(itemID)
		if err != nil {
			showError(w, "Routine item not found")
			return
		}

		if item.OrderIndex == 0 {
			showError(w, "Cannot move the first routine item up")
			return
		}

		err = db.MoveRoutineItem(item, true) // true for moving up
		if err != nil {
			showError(w, "Failed to move routine item up: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", item.RoutineID))
	}
}

func postRoutineItemsDown(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		itemID, err := g.GetIDFromPath(r)
		if err != nil {
			showError(w, "Invalid routine item ID: "+err.Error())
			return
		}

		item, err := db.GetRoutineItemByID(itemID)
		if err != nil {
			showError(w, "Routine item not found")
			return
		}

		if item.OrderIndex == len(item.Routine.RoutineItems)-1 {
			showError(w, "Cannot move the last routine item down")
			return
		}

		err = db.MoveRoutineItem(item, false) // false for moving down
		if err != nil {
			showError(w, "Failed to move routine item down: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", item.RoutineID))
	}
}
