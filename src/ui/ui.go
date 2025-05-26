package ui

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/birabittoh/go-lift/src/database"
)

const (
	base      = "base"
	templates = "templates"
	ps        = string(os.PathSeparator)

	basePath        = "templates" + ps + base + ".gohtml"
	errorPath       = "templates" + ps + "error.gohtml"
	homePath        = "templates" + ps + "home.gohtml"
	exercisesPath   = "templates" + ps + "exercises.gohtml"
	exercisePath    = "templates" + ps + "exercise.gohtml"
	routinesPath    = "templates" + ps + "routines.gohtml"
	workoutsPath    = "templates" + ps + "workouts.gohtml"
	profilePath     = "templates" + ps + "profile.gohtml"
	profileEditPath = "templates" + ps + "profile_edit.gohtml"
)

var (
	tmpl    map[string]*template.Template
	funcMap = template.FuncMap{
		"capitalize":      capitalize,
		"coalesce":        coalesce,
		"formatBirthDate": formatBirthDate,
	}
)

type PageData struct {
	Page      string
	Exercises []database.Exercise
	Users     []database.User
	Message   string
}

func parseTemplate(path string) *template.Template {
	return template.Must(template.New(base).Funcs(funcMap).ParseFiles(path, basePath))
}

func executeTemplateSafe(w http.ResponseWriter, t string, data any) {
	var buf bytes.Buffer
	if err := tmpl[t].ExecuteTemplate(&buf, base, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func showError(w http.ResponseWriter, message string) {
	pageData := &PageData{
		Page:    "error",
		Message: message,
	}

	executeTemplateSafe(w, errorPath, pageData)
}

func getHome(w http.ResponseWriter, r *http.Request) {
	executeTemplateSafe(w, homePath, &PageData{Page: "home"})
}

func getExercises(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exercises, err := db.GetExercises()
		if err != nil {
			showError(w, "Failed to retrieve exercises: "+err.Error())
			return
		}

		pageData := &PageData{
			Page:      "exercises",
			Exercises: exercises,
		}

		executeTemplateSafe(w, exercisesPath, pageData)
	}
}

func getExercise(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			showError(w, "Exercise ID is required")
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			showError(w, "Invalid exercise ID: "+err.Error())
			return
		}

		exercise, err := db.GetExerciseByID(uint(id))
		if err != nil {
			showError(w, "Failed to retrieve exercise: "+err.Error())
			return
		}

		pageData := &PageData{
			Page:      "exercises",
			Exercises: []database.Exercise{*exercise},
		}

		executeTemplateSafe(w, exercisePath, pageData)
	}
}

func getRoutines(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		executeTemplateSafe(w, routinesPath, &PageData{Page: "routines"})
	}
}

func getWorkouts(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		executeTemplateSafe(w, workoutsPath, &PageData{Page: "workouts"})
	}
}

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

func InitServeMux(s *http.ServeMux, db *database.Database) {

	tmpl = make(map[string]*template.Template)

	tmpl[errorPath] = parseTemplate(errorPath)
	tmpl[homePath] = parseTemplate(homePath)
	tmpl[exercisesPath] = parseTemplate(exercisesPath)
	tmpl[exercisePath] = parseTemplate(exercisePath)
	tmpl[routinesPath] = parseTemplate(routinesPath)
	tmpl[workoutsPath] = parseTemplate(workoutsPath)
	tmpl[profilePath] = parseTemplate(profilePath)
	tmpl[profileEditPath] = parseTemplate(profileEditPath)

	s.HandleFunc("GET /", getHome)
	s.HandleFunc("GET /exercises", getExercises(db))
	s.HandleFunc("GET /exercises/{id}", getExercise(db))
	s.HandleFunc("GET /routines", getRoutines(db))
	s.HandleFunc("GET /workouts", getWorkouts(db))
	s.HandleFunc("GET /profile", getProfile(db))
	s.HandleFunc("GET /profile/edit", getProfileEdit(db))
	s.HandleFunc("POST /profile/edit", postProfileEdit(db))

	s.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	})
}
