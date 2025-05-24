package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/birabittoh/go-lift/src/database"
	"gorm.io/gorm"
)

// User handlers
func getUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromPath(r)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		var user database.User
		if err := db.First(&user, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				jsonError(w, http.StatusNotFound, "User not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Database error")
			return
		}

		jsonResponse(w, http.StatusOK, user)
	}
}

func updateUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromPath(r)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		var user database.User
		if err := db.First(&user, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				jsonError(w, http.StatusNotFound, "User not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Database error")
			return
		}

		var updateData database.User
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// Update specific fields
		user.Name = updateData.Name
		user.IsFemale = updateData.IsFemale
		user.Height = updateData.Height
		user.Weight = updateData.Weight
		user.BirthDate = updateData.BirthDate

		if err := db.Save(&user).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to update user")
			return
		}

		jsonResponse(w, http.StatusOK, user)
	}
}

// Exercise handlers (read-only)
func getExercisesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var exercises []database.Exercise
		query := db.Preload("PrimaryMuscles").Preload("SecondaryMuscles")

		// Optional filtering
		if category := r.URL.Query().Get("category"); category != "" {
			query = query.Where("category = ?", category)
		}
		if level := r.URL.Query().Get("level"); level != "" {
			query = query.Where("level = ?", level)
		}

		if err := query.Find(&exercises).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Database error")
			return
		}

		jsonResponse(w, http.StatusOK, exercises)
	}
}

func getExerciseHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromPath(r)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid exercise ID")
			return
		}

		var exercise database.Exercise
		if err := db.Preload("PrimaryMuscles").Preload("SecondaryMuscles").First(&exercise, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				jsonError(w, http.StatusNotFound, "Exercise not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Database error")
			return
		}

		jsonResponse(w, http.StatusOK, exercise)
	}
}

// Muscle handlers (read-only)
func getMusclesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var muscles []database.Muscle
		if err := db.Find(&muscles).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Database error")
			return
		}

		jsonResponse(w, http.StatusOK, muscles)
	}
}

// Routine handlers
func getRoutinesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var routines []database.Routine
		if err := db.Find(&routines).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Database error")
			return
		}

		jsonResponse(w, http.StatusOK, routines)
	}
}

func getRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromPath(r)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid routine ID")
			return
		}

		var routine database.Routine
		if err := db.Preload("Items.ExerciseItems.Exercise.PrimaryMuscles").
			Preload("Items.ExerciseItems.Exercise.SecondaryMuscles").
			Preload("Items.ExerciseItems.Sets").
			Order("Items.order_index, Items.ExerciseItems.order_index, Items.ExerciseItems.Sets.order_index").
			First(&routine, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				jsonError(w, http.StatusNotFound, "Routine not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Database error")
			return
		}

		jsonResponse(w, http.StatusOK, routine)
	}
}

func createRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var routine database.Routine
		if err := json.NewDecoder(r.Body).Decode(&routine); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if err := db.Create(&routine).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to create routine")
			return
		}

		jsonResponse(w, http.StatusCreated, routine)
	}
}

func updateRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromPath(r)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid routine ID")
			return
		}

		var routine database.Routine
		if err := json.NewDecoder(r.Body).Decode(&routine); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		routine.ID = uint(id)
		if err := db.Save(&routine).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to update routine")
			return
		}

		jsonResponse(w, http.StatusOK, routine)
	}
}

func deleteRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromPath(r)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid routine ID")
			return
		}

		if err := db.Delete(&database.Routine{}, id).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to delete routine")
			return
		}

		jsonResponse(w, http.StatusOK, map[string]string{"message": "Routine deleted"})
	}
}

// Record routine handlers (workout sessions)
func getRecordRoutinesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var records []database.RecordRoutine
		query := db.Preload("Routine").Order("created_at DESC")

		// Optional limit for recent workouts
		if limit := r.URL.Query().Get("limit"); limit != "" {
			if l, err := strconv.Atoi(limit); err == nil {
				query = query.Limit(l)
			}
		}

		if err := query.Find(&records).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Database error")
			return
		}

		jsonResponse(w, http.StatusOK, records)
	}
}

func getRecordRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromPath(r)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid record ID")
			return
		}

		var record database.RecordRoutine
		if err := db.Preload("Routine").
			Preload("RecordItems.RoutineItem").
			Preload("RecordItems.RecordExerciseItems.ExerciseItem.Exercise.PrimaryMuscles").
			Preload("RecordItems.RecordExerciseItems.ExerciseItem.Exercise.SecondaryMuscles").
			Preload("RecordItems.RecordExerciseItems.RecordSets.Set").
			Order("RecordItems.order_index, RecordItems.RecordExerciseItems.order_index, RecordItems.RecordExerciseItems.RecordSets.order_index").
			First(&record, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				jsonError(w, http.StatusNotFound, "Record not found")
				return
			}
			jsonError(w, http.StatusInternalServerError, "Database error")
			return
		}

		jsonResponse(w, http.StatusOK, record)
	}
}

func createRecordRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var record database.RecordRoutine
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if err := db.Create(&record).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to create record")
			return
		}

		jsonResponse(w, http.StatusCreated, record)
	}
}

func updateRecordRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromPath(r)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid record ID")
			return
		}

		var record database.RecordRoutine
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		record.ID = uint(id)
		if err := db.Save(&record).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to update record")
			return
		}

		jsonResponse(w, http.StatusOK, record)
	}
}

func deleteRecordRoutineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromPath(r)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid record ID")
			return
		}

		if err := db.Delete(&database.RecordRoutine{}, id).Error; err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to delete record")
			return
		}

		jsonResponse(w, http.StatusOK, map[string]string{"message": "Record deleted"})
	}
}

// Stats handler
func getStatsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := WorkoutStats{}

		// Total workouts
		db.Model(&database.RecordRoutine{}).Count(&stats.TotalWorkouts)

		// Total minutes (sum of all workout durations)
		var totalSeconds uint
		db.Model(&database.RecordRoutine{}).Select("COALESCE(SUM(duration), 0)").Scan(&totalSeconds)
		stats.TotalMinutes = int(totalSeconds / 60)

		// Total exercises completed
		db.Model(&database.RecordExerciseItem{}).Count(&stats.TotalExercises)

		// Most frequent exercise
		var exerciseStats struct {
			ExerciseName string `json:"exercise_name"`
			Count        int    `json:"count"`
		}
		db.Raw(`
			SELECT e.name as exercise_name, COUNT(*) as count
			FROM record_exercise_items rei
			JOIN exercise_items ei ON rei.exercise_item_id = ei.id
			JOIN exercises e ON ei.exercise_id = e.id
			GROUP BY e.id, e.name
			ORDER BY count DESC
			LIMIT 1
		`).Scan(&exerciseStats)

		if exerciseStats.Count > 0 {
			stats.MostFrequentExercise = &struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			}{
				Name:  exerciseStats.ExerciseName,
				Count: exerciseStats.Count,
			}
		}

		// Most frequent routine
		var routineStats struct {
			RoutineName string `json:"routine_name"`
			Count       int    `json:"count"`
		}
		db.Raw(`
			SELECT r.name as routine_name, COUNT(*) as count
			FROM record_routines rr
			JOIN routines r ON rr.routine_id = r.id
			GROUP BY r.id, r.name
			ORDER BY count DESC
			LIMIT 1
		`).Scan(&routineStats)

		if routineStats.Count > 0 {
			stats.MostFrequentRoutine = &struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			}{
				Name:  routineStats.RoutineName,
				Count: routineStats.Count,
			}
		}

		// Recent workouts (last 5)
		db.Preload("Routine").
			Order("created_at DESC").
			Limit(5).
			Find(&stats.RecentWorkouts)

		jsonResponse(w, http.StatusOK, stats)
	}
}
