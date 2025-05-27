package ui

import "net/http"

func getHome(w http.ResponseWriter, r *http.Request) {
	executeTemplateSafe(w, homePath, &PageData{Page: "home"})
}
