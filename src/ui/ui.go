package ui

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

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
	Routines  []database.Routine
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
