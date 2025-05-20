package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
	"gorm.io/gorm"
)

// Routines handlers
func getRoutinesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var routines []database.Routine
		result := db.Preload("Exercises").Preload("Supersets").Preload("Supersets.Sets").Find(&routines)
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
		result := db.Preload("Exercises").Preload("Supersets").Preload("Supersets.Sets").First(&routine, id)
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
