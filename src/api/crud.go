package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
	"gorm.io/gorm"
)

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

// User handlers
func getUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var user database.User
		if err := db.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				jsonError(w, http.StatusNotFound, "User not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, user)
	}
}

func updateUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var user database.User
		if err := db.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				jsonError(w, http.StatusNotFound, "User not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}
		if err := db.Save(&user).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to update user: "+err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, user)
	}
}

// Routines handlers
func getRoutinesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var routines []database.Routine
		result := db.
			Preload("RoutineItems").
			Preload("RoutineItems.Exercises").
			Preload("RoutineItems.Supersets").
			Preload("RoutineItems.Supersets.PrimaryExercise").
			Preload("RoutineItems.Supersets.SecondaryExercise").
			Find(&routines)
		if result.Error != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to fetch routines: "+result.Error.Error())
			return
		}
		jsonResponse(w, http.StatusOK, routines)
	}
}

func getRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var routine database.Routine
		result := db.
			Preload("RoutineItems.Exercises").
			Preload("RoutineItems.Supersets").
			Preload("RoutineItems.Supersets.PrimaryExercise").
			Preload("RoutineItems.Supersets.SecondaryExercise").
			First(&routine, id)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				jsonError(w, http.StatusNotFound, "Routine not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Failed to fetch routine: "+result.Error.Error())
			return
		}
		jsonResponse(w, http.StatusOK, routine)
	}
}

func createRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var routine database.Routine
		if err := json.NewDecoder(r.Body).Decode(&routine); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}

		if err := db.Create(&routine).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to create routine: "+err.Error())
			return
		}

		// Reload with associations
		db.Preload("Exercises").Preload("Supersets").Preload("Supersets.Sets").First(&routine, routine.ID)

		jsonResponse(w, http.StatusCreated, routine)
	}
}

func updateRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var routine database.Routine

		// Check if exists
		if err := db.First(&routine, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				jsonError(w, http.StatusNotFound, "Routine not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}

		// Parse update data
		if err := json.NewDecoder(r.Body).Decode(&routine); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}

		// Save with associations
		if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&routine).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to update routine: "+err.Error())
			return
		}

		// Reload complete data
		db.Preload("Exercises").Preload("Supersets").Preload("Supersets.Sets").First(&routine, id)
		jsonResponse(w, http.StatusOK, routine)
	}
}

func deleteRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := db.Delete(&database.Routine{}, id).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to delete routine: "+err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"message": "Routine deleted successfully"})
	}
}

// Exercises handlers
func getExercisesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var exercises []database.Exercise
		if err := db.Find(&exercises).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to fetch exercises: "+err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, exercises)
	}
}

func getExerciseHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var exercise database.Exercise
		if err := db.First(&exercise, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				jsonError(w, http.StatusNotFound, "Exercise not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Failed to fetch exercise: "+err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, exercise)
	}
}

func createExerciseHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var exercise database.Exercise
		if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}

		if err := db.Create(&exercise).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to create exercise: "+err.Error())
			return
		}

		jsonResponse(w, http.StatusCreated, exercise)
	}
}

func upsertExercisesHandler(dbStruct *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := dbStruct.UpdateExercises(); err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to update exercises: "+err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, map[string]string{"message": "Exercises updated successfully"})
	}
}

func updateExerciseHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		// Verify exercise exists
		var exercise database.Exercise
		if err := db.First(&exercise, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				jsonError(w, http.StatusNotFound, "Exercise not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}

		// Parse update data
		if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}

		if err := db.Save(&exercise).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to update exercise: "+err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, exercise)
	}
}

func deleteExerciseHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := db.Delete(&database.Exercise{}, id).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to delete exercise: "+err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"message": "Exercise deleted successfully"})
	}
}

// RecordRoutines handlers
func getRecordRoutinesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var records []database.RecordRoutine
		result := db.Preload("RecordExercises").Preload("Routine.Exercises").Preload("Routine.Supersets").Find(&records)
		if result.Error != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to fetch record routines: "+result.Error.Error())
			return
		}
		jsonResponse(w, http.StatusOK, records)
	}
}

func getRecordRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var record database.RecordRoutine
		result := db.Preload("Routine").Preload("Routine.Exercises").Preload("Routine.Supersets").First(&record, id)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				jsonError(w, http.StatusNotFound, "Record routine not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Failed to fetch record routine: "+result.Error.Error())
			return
		}
		jsonResponse(w, http.StatusOK, record)
	}
}

func createRecordRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var record database.RecordRoutine
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}

		if err := db.Create(&record).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to create record routine: "+err.Error())
			return
		}

		// Reload with associations
		db.Preload("Routine").Preload("Routine.Exercises").Preload("Routine.Supersets").First(&record, record.ID)

		jsonResponse(w, http.StatusCreated, record)
	}
}

func updateRecordRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var record database.RecordRoutine

		// Check if exists
		if err := db.First(&record, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				jsonError(w, http.StatusNotFound, "Record routine not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}

		// Parse update data
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}

		// Save with associations
		if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&record).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to update record routine: "+err.Error())
			return
		}

		// Reload complete data
		db.Preload("Routine").Preload("Routine.Exercises").Preload("Routine.Supersets").First(&record, id)
		jsonResponse(w, http.StatusOK, record)
	}
}

func deleteRecordRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := db.Delete(&database.RecordRoutine{}, id).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to delete record routine: "+err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"message": "Record routine deleted successfully"})
	}
}

func getStatsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := WorkoutStats{}

		// Get total workouts
		if err := db.Model(&database.RecordRoutine{}).Count(&stats.TotalWorkouts).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to count workouts: "+err.Error())
			return
		}

		// Get total minutes
		stats.TotalMinutes = 0

		// Get total exercises
		if err := db.Model(&database.RecordExercise{}).Count(&stats.TotalExercises).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to count exercises: "+err.Error())
			return
		}

		// Get most frequent exercise
		var mostFrequentExercise struct {
			Name  string `gorm:"column:name"`
			Count int    `gorm:"column:count"`
		}
		exerciseQuery := db.Model(&database.RecordExercise{}).
			Select("exercises.name, COUNT(*) as count").
			Joins("JOIN exercises ON record_exercises.exercise_id = exercises.id").
			Group("exercises.name").
			Order("count DESC").
			Limit(1)

		if err := exerciseQuery.Scan(&mostFrequentExercise).Error; err == nil && mostFrequentExercise.Name != "" {
			stats.MostFrequentExercise = &struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			}{
				Name:  mostFrequentExercise.Name,
				Count: mostFrequentExercise.Count,
			}
		}

		// Get most frequent routine
		var mostFrequentRoutine struct {
			Name  string `gorm:"column:name"`
			Count int    `gorm:"column:count"`
		}
		routineQuery := db.Model(&database.RecordRoutine{}).
			Select("routines.name, COUNT(*) as count").
			Joins("JOIN routines ON record_routines.routine_id = routines.id").
			Group("routines.name").
			Order("count DESC").
			Limit(1)

		if err := routineQuery.Scan(&mostFrequentRoutine).Error; err == nil && mostFrequentRoutine.Name != "" {
			stats.MostFrequentRoutine = &struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			}{
				Name:  mostFrequentRoutine.Name,
				Count: mostFrequentRoutine.Count,
			}
		}

		// Get recent workouts (last 5)
		if err := db.
			Preload("RecordRoutineItems").
			Preload("Routine").
			Order("created_at DESC").
			Limit(5).
			Find(&stats.RecentWorkouts).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to fetch recent workouts: "+err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, stats)
	}
}
