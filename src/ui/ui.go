package ui

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
)

const (
	base      = "base"
	templates = "templates"
	ps        = string(os.PathSeparator)

	basePath      = "templates" + ps + base + ".gohtml"
	homePath      = "templates" + ps + "home.gohtml"
	exercisesPath = "templates" + ps + "exercises.gohtml"
	routinesPath  = "templates" + ps + "routines.gohtml"
	workoutsPath  = "templates" + ps + "workouts.gohtml"
	profilePath   = "templates" + ps + "profile.gohtml"
)

var (
	tmpl    map[string]*template.Template
	funcMap = template.FuncMap{
		"capitalize": capitalize,
	}
)

type PageData struct {
	Page string
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

func getHome(w http.ResponseWriter, r *http.Request) {
	executeTemplateSafe(w, homePath, &PageData{Page: "home"})
}

func getExercises(w http.ResponseWriter, r *http.Request) {
	executeTemplateSafe(w, exercisesPath, &PageData{Page: "exercises"})
}

func getRoutines(w http.ResponseWriter, r *http.Request) {
	executeTemplateSafe(w, routinesPath, &PageData{Page: "routines"})
}

func getWorkouts(w http.ResponseWriter, r *http.Request) {
	executeTemplateSafe(w, workoutsPath, &PageData{Page: "workouts"})
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	executeTemplateSafe(w, profilePath, &PageData{Page: "profile"})
}

func InitServeMux(s *http.ServeMux) {
	tmpl = make(map[string]*template.Template)

	tmpl[homePath] = parseTemplate(homePath)
	tmpl[exercisesPath] = parseTemplate(exercisesPath)
	tmpl[routinesPath] = parseTemplate(routinesPath)
	tmpl[workoutsPath] = parseTemplate(workoutsPath)
	tmpl[profilePath] = parseTemplate(profilePath)

	s.HandleFunc("GET /", getHome)
	s.HandleFunc("GET /exercises", getExercises)
	s.HandleFunc("GET /routines", getRoutines)
	s.HandleFunc("GET /workouts", getWorkouts)
	s.HandleFunc("GET /profile", getProfile)

	s.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	})
}
