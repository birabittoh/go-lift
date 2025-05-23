package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/birabittoh/go-lift/src/database"
)

const uiDir = "ui/dist"

var fileServer = http.FileServer(http.Dir(uiDir))

func GetServeMux(dbStruct *database.Database) *http.ServeMux {
	mux := http.NewServeMux()
	db := dbStruct.DB

	mux.HandleFunc("GET /authelia/api/user/info", mockAutheliaHandler)

	mux.HandleFunc("GET /api/ping", pingHandler)
	mux.HandleFunc("GET /api/connection", connectionHandler(db))

	// Profile routes
	mux.HandleFunc("GET /api/users/{id}", getUserHandler(db))
	mux.HandleFunc("PUT /api/users/{id}", updateUserHandler(db))

	// Routines routes
	mux.HandleFunc("GET /api/routines", getRoutinesHandler(db))
	mux.HandleFunc("GET /api/routines/{id}", getRoutineHandler(db))
	mux.HandleFunc("POST /api/routines", createRoutineHandler(db))
	mux.HandleFunc("PUT /api/routines/{id}", updateRoutineHandler(db))
	mux.HandleFunc("DELETE /api/routines/{id}", deleteRoutineHandler(db))

	// Exercises routes
	mux.HandleFunc("GET /api/exercises", getExercisesHandler(db))
	mux.HandleFunc("GET /api/exercises/{id}", getExerciseHandler(db))
	mux.HandleFunc("POST /api/exercises", createExerciseHandler(db))
	mux.HandleFunc("POST /api/exercises/update", upsertExercisesHandler(dbStruct))
	mux.HandleFunc("PUT /api/exercises/{id}", updateExerciseHandler(db))
	mux.HandleFunc("DELETE /api/exercises/{id}", deleteExerciseHandler(db))

	// RecordRoutines routes
	mux.HandleFunc("GET /api/recordroutines", getRecordRoutinesHandler(db))
	mux.HandleFunc("GET /api/recordroutines/{id}", getRecordRoutineHandler(db))
	mux.HandleFunc("POST /api/recordroutines", createRecordRoutineHandler(db))
	mux.HandleFunc("PUT /api/recordroutines/{id}", updateRecordRoutineHandler(db))
	mux.HandleFunc("DELETE /api/recordroutines/{id}", deleteRecordRoutineHandler(db))

	// Stats routes
	mux.HandleFunc("GET /api/stats", getStatsHandler(db))

	// Static UI route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the file exists at the requested path
		requestedFile := filepath.Join(uiDir, r.URL.Path)
		_, err := os.Stat(requestedFile)

		// If file exists or it's the root, serve it directly
		if err == nil || r.URL.Path == "/" {
			fileServer.ServeHTTP(w, r)
			return
		}

		// For file not found, serve index.html for SPA routing
		http.ServeFile(w, r, filepath.Join(uiDir, "index.html"))
	})

	return mux
}
