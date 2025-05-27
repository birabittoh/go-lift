package ui

import (
	"net/http"
	"strconv"
	"time"

	"github.com/birabittoh/go-lift/src/database"
)

func getProfile(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := db.GetUserByID(1)
		if err != nil {
			showError(w, "Failed to retrieve profile: "+err.Error())
			return
		}

		pageData := &PageData{
			Page:  "profile",
			Users: []database.User{*user},
		}
		executeTemplateSafe(w, profilePath, pageData)
	}
}

func getProfileEdit(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := db.GetUserByID(1)
		if err != nil {
			showError(w, "Failed to retrieve profile: "+err.Error())
			return
		}

		pageData := &PageData{
			Page:  "profile",
			Users: []database.User{*user},
		}
		executeTemplateSafe(w, profileEditPath, pageData)
	}
}

func postProfileEdit(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			showError(w, "Failed to parse form: "+err.Error())
			return
		}

		// Get the current user
		user, err := db.GetUserByID(1)
		if err != nil {
			showError(w, "Failed to retrieve user: "+err.Error())
			return
		}

		// Update user fields from form
		user.Name = r.FormValue("name")
		user.IsFemale = r.FormValue("isFemale") == "true"

		heightValue := r.FormValue("height")
		if heightValue == "" {
			user.Height = nil
		} else {
			height, err := strconv.ParseFloat(heightValue, 64)
			if err != nil {
				showError(w, "Invalid height value: "+err.Error())
				return
			}
			user.Height = &height
		}

		weightValue := r.FormValue("weight")
		if weightValue == "" {
			user.Weight = nil
		} else {
			weight, err := strconv.ParseFloat(weightValue, 64)
			if err != nil {
				showError(w, "Invalid weight value: "+err.Error())
				return
			}
			user.Weight = &weight
		}

		birthDateValue, err := time.Parse("2006-01-02", r.FormValue("birthDate"))
		if err != nil {
			user.BirthDate = nil
		} else {
			user.BirthDate = &birthDateValue
		}

		// Save the updated user
		err = db.UpdateUser(user)
		if err != nil {
			showError(w, "Failed to update profile: "+err.Error())
			return
		}

		// Redirect to profile page
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
