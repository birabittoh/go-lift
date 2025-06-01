package ui

import (
	"fmt"
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
	g "github.com/birabittoh/go-lift/src/globals"
)

func postAddRecordRoutine(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cw := db.GetCurrentWorkout()
		if cw != nil {
			showError(w, "A workout is already in progress")
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

		rr, err := db.NewRecordRoutine(routine)
		if err != nil {
			showError(w, "Failed to start routine: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/record-routines/%d", rr.ID))
	}
}

func postRecordRoutinesDelete(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, err := g.GetIDFromPath(r)
		if err != nil {
			showError(w, "Invalid record routine ID: "+err.Error())
			return
		}

		recordRoutine, err := db.GetRecordRoutineByID(id)
		if err != nil {
			showError(w, "Record routine not found")
			return
		}

		err = db.DeleteRecordRoutine(recordRoutine)
		if err != nil {
			showError(w, "Failed to delete record routine: "+err.Error())
			return
		}

		// get "page" query parameter
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "home"
		}

		redirect(w, r, "/"+page)
	}
}
