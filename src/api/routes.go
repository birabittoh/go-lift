package api

import (
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
	"github.com/birabittoh/go-lift/src/ui"
)

func GetServeMux(db *database.Database) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /authelia/api/user/info", mockAutheliaHandler)

	mux.HandleFunc("GET /api/ping", pingHandler)
	mux.HandleFunc("GET /api/connection", connectionHandler(db))

	// Profile routes
	mux.HandleFunc("GET /api/users/{id}", getUserHandler(db))
	mux.HandleFunc("PUT /api/users/{id}", updateUserHandler(db))

	// Exercises routes (read-only)
	mux.HandleFunc("GET /api/exercises", getExercisesHandler(db))
	mux.HandleFunc("GET /api/exercises/{id}", getExerciseHandler(db))

	// Routines routes
	mux.HandleFunc("GET /api/routines", getRoutinesHandler(db))
	mux.HandleFunc("GET /api/routines/{id}", getRoutineHandler(db))
	mux.HandleFunc("POST /api/routines", createRoutineHandler(db))
	mux.HandleFunc("PUT /api/routines/{id}", updateRoutineHandler(db))
	mux.HandleFunc("DELETE /api/routines/{id}", deleteRoutineHandler(db))

	// Record routines routes (workout sessions)
	mux.HandleFunc("GET /api/records", getRecordRoutinesHandler(db))
	mux.HandleFunc("GET /api/records/{id}", getRecordRoutineHandler(db))
	mux.HandleFunc("POST /api/records", createRecordRoutineHandler(db))
	mux.HandleFunc("PUT /api/records/{id}", updateRecordRoutineHandler(db))
	mux.HandleFunc("DELETE /api/records/{id}", deleteRecordRoutineHandler(db))

	// Stats routes
	mux.HandleFunc("GET /api/stats", getStatsHandler(db))

	ui.InitServeMux(mux, db)

	return mux
}

type WorkoutStats struct {
	TotalWorkouts        int64 `json:"totalWorkouts"`
	TotalMinutes         int   `json:"totalMinutes"`
	TotalExercises       int64 `json:"totalExercises"`
	MostFrequentExercise *struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	} `json:"mostFrequentExercise,omitempty"`
	MostFrequentRoutine *struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	} `json:"mostFrequentRoutine,omitempty"`
	RecentWorkouts []database.RecordRoutine `json:"recentWorkouts"`
}
