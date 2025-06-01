package ui

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	"github.com/birabittoh/go-lift/src/database"
	g "github.com/birabittoh/go-lift/src/globals"
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
	routinePath     = "templates" + ps + "routine.gohtml"
	workoutsPath    = "templates" + ps + "workouts.gohtml"
	profilePath     = "templates" + ps + "profile.gohtml"
	profileEditPath = "templates" + ps + "profile_edit.gohtml"
)

var (
	tmpl    map[string]*template.Template
	funcMap = template.FuncMap{
		"capitalize":      g.Capitalize,
		"coalesce":        coalesce,
		"formatBirthDate": formatBirthDate,
		"formatDay":       formatDay,
		"isChecked":       isChecked,
	}
)

type PageData struct {
	Page      string
	Days      []database.Day
	Exercises []database.Exercise
	Routines  []database.Routine
	Users     []database.User
	Message   string
	ID        uint
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
	tmpl[routinePath] = parseTemplate(routinePath)
	tmpl[workoutsPath] = parseTemplate(workoutsPath)
	tmpl[profilePath] = parseTemplate(profilePath)
	tmpl[profileEditPath] = parseTemplate(profileEditPath)

	s.HandleFunc("GET /", getHome)                                    // home page
	s.HandleFunc("GET /exercises/{id}", getExercises(db))             // select exercise for routine item id
	s.HandleFunc("GET /exercises/{id}/{exerciseId}", getExercise(db)) // confirm exercise for routine item id
	s.HandleFunc("GET /routines", getRoutines(db))                    // list all routines
	s.HandleFunc("GET /routines/{id}", getRoutine(db))                // edit routine
	s.HandleFunc("GET /workouts", getWorkouts(db))                    // list all workouts
	s.HandleFunc("GET /profile", getProfile(db))                      // user profile
	s.HandleFunc("GET /profile/edit", getProfileEdit(db))             // edit user profile

	s.HandleFunc("POST /exercises/{id}/{exerciseId}", postAddExercise(db))       // add exercise item to routine item
	s.HandleFunc("POST /routines/new", postRoutineNew(db))                       // add new routine
	s.HandleFunc("POST /routines/{id}", postRoutine(db))                         // edit routine (name, description)
	s.HandleFunc("POST /routines/{id}/delete", postRoutineDelete(db))            // delete routine
	s.HandleFunc("POST /routines/{id}/new", postRoutineItemNew(db))              // add new routine item to routine
	s.HandleFunc("POST /exercise-items/{id}/delete", postExerciseItemDelete(db)) // delete exercise item
	s.HandleFunc("POST /exercise-items/{id}", postExerciseItem(db))              // edit exercise item (restTime, sets)
	s.HandleFunc("POST /exercise-items/{id}/new", postSetNew(db))                // add new set to exercise item
	s.HandleFunc("POST /sets/{id}/delete", postSetDelete(db))                    // delete set
	s.HandleFunc("POST /profile/edit", postProfileEdit(db))                      // edit user profile

	s.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	})
}
