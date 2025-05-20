package api

import (
	"net/http"

	"gorm.io/gorm"
)

func GetServeMux(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /authelia/api/user/info", mockAutheliaHandler)

	mux.HandleFunc("GET /api/ping", pingHandler)
	mux.HandleFunc("GET /api/connection", connectionHandler(db))

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
	mux.HandleFunc("PUT /api/exercises/{id}", updateExerciseHandler(db))
	mux.HandleFunc("DELETE /api/exercises/{id}", deleteExerciseHandler(db))

	// RecordRoutines routes
	mux.HandleFunc("GET /api/recordroutines", getRecordRoutinesHandler(db))
	mux.HandleFunc("GET /api/recordroutines/{id}", getRecordRoutineHandler(db))
	mux.HandleFunc("POST /api/recordroutines", createRecordRoutineHandler(db))
	mux.HandleFunc("PUT /api/recordroutines/{id}", updateRecordRoutineHandler(db))
	mux.HandleFunc("DELETE /api/recordroutines/{id}", deleteRecordRoutineHandler(db))

	return mux
}
