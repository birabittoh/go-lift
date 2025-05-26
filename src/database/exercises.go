package database

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ImportedExercise represents the JSON structure from the input file
type importedExercise struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Level            string   `json:"level"`
	Category         string   `json:"category"`
	Force            *string  `json:"force"`
	Mechanic         *string  `json:"mechanic"`
	Equipment        *string  `json:"equipment"`
	PrimaryMuscles   []string `json:"primaryMuscles"`
	SecondaryMuscles []string `json:"secondaryMuscles"`
	Instructions     []string `json:"instructions"`
	Images           []string `json:"images"`
}

// upsertExercise creates or updates a single exercise with all its related data
func (db *Database) upsertExercise(importedExercise importedExercise) (didSave, isUpdate bool, err error) {
	// First, try to find existing exercise by name
	var existingExercise Exercise
	result := db.Where("id = ?", importedExercise.ID).Preload("PrimaryMuscles").Preload("SecondaryMuscles").First(&existingExercise)

	// Create new exercise with basic info
	exercise := Exercise{
		ID:       importedExercise.ID,
		Name:     importedExercise.Name,
		Level:    importedExercise.Level,
		Category: importedExercise.Category,
	}

	if importedExercise.Force != nil && *importedExercise.Force != "" {
		exercise.Force = importedExercise.Force
	}
	if importedExercise.Mechanic != nil && *importedExercise.Mechanic != "" {
		exercise.Mechanic = importedExercise.Mechanic
	}
	if importedExercise.Equipment != nil && *importedExercise.Equipment != "" {
		exercise.Equipment = importedExercise.Equipment
	}
	if len(importedExercise.PrimaryMuscles) > 0 {
		primaryMuscles := strings.Join(importedExercise.PrimaryMuscles, ", ")
		exercise.PrimaryMuscles = &primaryMuscles
	}
	if len(importedExercise.SecondaryMuscles) > 0 {
		secondaryMuscles := strings.Join(importedExercise.SecondaryMuscles, ", ")
		exercise.SecondaryMuscles = &secondaryMuscles
	}

	if len(importedExercise.Instructions) > 0 {
		// Filter out empty instructions
		var filteredInstructions []string
		for _, instruction := range importedExercise.Instructions {
			clean := strings.TrimSpace(instruction)
			if clean != "" {
				filteredInstructions = append(filteredInstructions, clean)
			}
		}
		instructions := strings.Join(filteredInstructions, "\n")
		exercise.InstructionsString = &instructions
	}

	var exerciseDataChanged bool

	if result.Error == nil {
		// Exercise exists, check if it needs updating
		isUpdate = true
		exercise.ID = existingExercise.ID
		exercise.CreatedAt = existingExercise.CreatedAt // Preserve creation time

		// Check if the exercise data has actually changed
		exerciseDataChanged = db.exerciseDataChanged(existingExercise, exercise)

		// Only update if something has changed
		if exerciseDataChanged {
			if err := db.Save(&exercise).Error; err != nil {
				return false, false, fmt.Errorf("failed to update exercise: %w", err)
			}
			didSave = true
		}
	} else {
		// Exercise doesn't exist, create it
		isUpdate = false
		exerciseDataChanged = true // New exercise, so data is "changed"

		if err := db.Create(&exercise).Error; err != nil {
			return false, false, fmt.Errorf("failed to create exercise: %w", err)
		}
		didSave = true
	}

	return
}

// exerciseDataChanged compares two exercises to see if core data has changed
func (db *Database) exerciseDataChanged(existing, new Exercise) bool {
	return existing.Level != new.Level ||
		existing.Category != new.Category ||
		!stringPointersEqual(existing.Force, new.Force) ||
		!stringPointersEqual(existing.Mechanic, new.Mechanic) ||
		!stringPointersEqual(existing.Equipment, new.Equipment) ||
		!stringPointersEqual(existing.InstructionsString, new.InstructionsString) ||
		existing.PrimaryMuscles != new.PrimaryMuscles ||
		existing.SecondaryMuscles != new.SecondaryMuscles
}

// Helper function to compare string pointers
func stringPointersEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// downloadExercises downloads exercises from the JSON URL
func downloadExercises() ([]importedExercise, error) {
	// Download exercises.json from the URL
	resp, err := http.Get(jsonURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download exercises.json: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download exercises.json: HTTP status %d", resp.StatusCode)
	}

	// Read the response body
	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var exercises []importedExercise
	if err := json.Unmarshal(fileData, &exercises); err != nil {
		return nil, fmt.Errorf("failed to parse exercises.json: %w", err)
	}

	return exercises, nil
}

const (
	dbDir       = "data"
	dbName      = "fitness.sqlite"
	baseURL     = "https://raw.githubusercontent.com/yuhonas/free-exercise-db/main/"
	jsonURL     = baseURL + "dist/exercises.json"
	imageFormat = baseURL + "exercises/%s/%d.jpg"
	imageAmount = 2
)

var lastUpdate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

func (e Exercise) GetImages() (images []string) {
	for i := range imageAmount {
		images = append(images, fmt.Sprintf(imageFormat, e.ID, i))
	}
	return
}

func (e Exercise) GetInstructions() []string {
	if e.InstructionsString == nil {
		return nil
	}

	return strings.Split(*e.InstructionsString, "\n")
}

func (e *Exercise) Fill() {
	e.Images = e.GetImages()
	e.Instructions = e.GetInstructions()
}

func (db *Database) UpdateExercises() (err error) {
	// Load exercises
	exercises, err := downloadExercises()
	if err != nil {
		log.Fatalf("Failed to load exercises: %v", err)
	}

	log.Printf("Successfully loaded %d exercises from JSON", len(exercises))

	var successCount, createCount, updateCount int

	// Import/update exercises
	for i, exercise := range exercises {
		didSave, isUpdate, err := db.upsertExercise(exercise)
		if err != nil {
			log.Printf("Failed to upsert exercise %d (%s): %v", i+1, exercise.Name, err)
			continue
		}

		successCount++
		if didSave {
			if isUpdate {
				updateCount++
			} else {
				createCount++
			}
		}
	}

	lastUpdate = time.Now()

	log.Printf("Update completed successfully! Processed %d out of %d exercises (%d created, %d updated)", successCount, len(exercises), createCount, updateCount)
	return
}

func (db *Database) GetLastUpdate() time.Time {
	return lastUpdate
}
