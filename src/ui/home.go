package ui

import (
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
)

func getHome(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageData, err := getPageData(db, "home")
		if err != nil {
			showError(w, "Failed to retrieve page data: "+err.Error())
			return
		}

		executeTemplateSafe(w, homePath, pageData)
	}
}
