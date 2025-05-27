package ui

import (
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
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
