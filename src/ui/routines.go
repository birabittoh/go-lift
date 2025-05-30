package ui

import (
	"fmt"
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
	g "github.com/birabittoh/go-lift/src/globals"
)

func getRoutines(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routines, err := db.GetRoutines()
		if err != nil {
			showError(w, "Failed to retrieve routines: "+err.Error())
			return
		}
		pageData := &PageData{
			Page:     "routines",
			Routines: routines,
		}
		executeTemplateSafe(w, routinesPath, pageData)
	}
}

func getRoutine(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		pageData := &PageData{
			Page:     "routines",
			Routines: []database.Routine{*routine},
			Days:     db.GetDays(),
		}
		executeTemplateSafe(w, routinePath, pageData)
	}
}

func postRoutineNew(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routine := &database.Routine{
			Name:        "New Routine",
			Description: "",
		}

		err := db.NewRoutine(routine)
		if err != nil {
			showError(w, "Failed to create new routine: "+err.Error())
			return
		}

		// Redirect to edit page
		redirect(w, r, fmt.Sprintf("/routines/%d", routine.ID))
	}
}

func postRoutineDelete(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := g.GetIDFromPath(r)
		if err != nil {
			showError(w, "Invalid routine ID: "+err.Error())
			return
		}

		err = db.DeleteRoutine(id)
		if err != nil {
			showError(w, "Failed to delete routine: "+err.Error())
			return
		}

		redirect(w, r, "/routines")
	}
}

func postRoutine(db *database.Database) http.HandlerFunc {
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

		routine.Name = r.FormValue("name")
		routine.Description = r.FormValue("description")

		weekDays := db.GetDays()
		var days []string
		for _, day := range weekDays {
			days = append(days, r.FormValue(day.Name))
		}

		routineDays := make([]database.Day, 0, len(days))
		for i, value := range days {
			if value == "on" {
				routineDays = append(routineDays, weekDays[i])
			}
		}

		err = db.UpdateRoutine(routine)
		if err != nil {
			showError(w, "Failed to update routine: "+err.Error())
			return
		}

		err = db.UpdateRoutineDays(routine, routineDays)
		if err != nil {
			showError(w, "Failed to update routine days: "+err.Error())
			return
		}

		redirect(w, r, fmt.Sprintf("/routines/%d", routine.ID))
	}
}
