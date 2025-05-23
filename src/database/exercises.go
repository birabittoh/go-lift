package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ImportedExercise represents the JSON structure from the input file
type importedExercise struct {
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
	ID               string   `json:"id"`
}

// upsertLookupEntities creates or updates all lookup table entries and returns lookup maps
func (db *Database) upsertLookupEntities(exercises []importedExercise) (map[string]uint, error) {
	// Collect unique values
	uniqueValues := make(map[string]bool)

	// Extract unique values from exercises
	for _, exercise := range exercises {
		for _, muscle := range exercise.PrimaryMuscles {
			uniqueValues[muscle] = true
		}
		for _, muscle := range exercise.SecondaryMuscles {
			uniqueValues[muscle] = true
		}
	}

	// Upsert lookup entities in database
	if err := db.upsertMuscles(uniqueValues); err != nil {
		return nil, err
	}

	// Build lookup map
	muscleLookupMap, err := db.buildMuscleLookupMap(uniqueValues)
	if err != nil {
		return nil, err
	}

	log.Printf("Upserted lookup entities: %d muscles", len(uniqueValues))

	return muscleLookupMap, nil
}

func (db *Database) upsertMuscles(muscles map[string]bool) error {
	for name := range muscles {
		muscle := Muscle{Name: name}
		if err := db.Where("name = ?", name).FirstOrCreate(&muscle).Error; err != nil {
			return fmt.Errorf("failed to upsert muscle %s: %w", name, err)
		}
	}
	return nil
}

// buildLookupMaps populates the lookup maps with IDs from the database
func (db *Database) buildMuscleLookupMap(uniqueValues map[string]bool) (map[string]uint, error) {
	// Build muscles map
	maps := make(map[string]uint)
	for name := range uniqueValues {
		var muscle Muscle
		if err := db.Where("name = ?", name).First(&muscle).Error; err != nil {
			return nil, fmt.Errorf("failed to find muscle %s: %w", name, err)
		}
		maps[name] = muscle.ID
	}

	return maps, nil
}

// upsertExercise creates or updates a single exercise with all its related data
func (db *Database) upsertExercise(importedExercise importedExercise, lookupMap map[string]uint) (didSave, isUpdate bool, err error) {
	// First, try to find existing exercise by name
	var existingExercise Exercise
	result := db.Where("name = ?", importedExercise.Name).Preload("PrimaryMuscles").Preload("SecondaryMuscles").First(&existingExercise)

	// Create new exercise with basic info
	exercise := Exercise{
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
		exercise.Instructions = &instructions
	}

	var exerciseID uint
	var exerciseDataChanged bool
	var muscleAssociationsChanged bool

	if result.Error == nil {
		// Exercise exists, check if it needs updating
		isUpdate = true
		exerciseID = existingExercise.ID
		exercise.ID = exerciseID
		exercise.CreatedAt = existingExercise.CreatedAt // Preserve creation time

		// Check if the exercise data has actually changed
		exerciseDataChanged = db.exerciseDataChanged(existingExercise, exercise)

		// Check if muscle associations have changed
		muscleAssociationsChanged = db.muscleAssociationsChanged(existingExercise, importedExercise, lookupMap)

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
		exerciseDataChanged = true       // New exercise, so data is "changed"
		muscleAssociationsChanged = true // New exercise, so associations are "changed"

		if err := db.Create(&exercise).Error; err != nil {
			return false, false, fmt.Errorf("failed to create exercise: %w", err)
		}
		exerciseID = exercise.ID
		didSave = true
	}

	// Only update muscle associations if they've changed
	if muscleAssociationsChanged {
		if err := db.updateMuscleAssociations(exerciseID, importedExercise, lookupMap); err != nil {
			return false, false, fmt.Errorf("failed to update muscle associations: %w", err)
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
		!stringPointersEqual(existing.Instructions, new.Instructions)
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

// muscleAssociationsChanged compares existing muscle associations with imported data
func (db *Database) muscleAssociationsChanged(existing Exercise, imported importedExercise, lookupMap map[string]uint) bool {
	// Convert existing muscle associations to sets for comparison
	existingPrimary := make(map[uint]bool)
	existingSecondary := make(map[uint]bool)

	for _, muscle := range existing.PrimaryMuscles {
		existingPrimary[muscle.ID] = true
	}
	for _, muscle := range existing.SecondaryMuscles {
		existingSecondary[muscle.ID] = true
	}

	// Convert imported muscle names to IDs and create sets
	importedPrimary := make(map[uint]bool)
	importedSecondary := make(map[uint]bool)

	for _, muscleName := range imported.PrimaryMuscles {
		if muscleID, ok := lookupMap[muscleName]; ok {
			importedPrimary[muscleID] = true
		}
	}
	for _, muscleName := range imported.SecondaryMuscles {
		if muscleID, ok := lookupMap[muscleName]; ok {
			importedSecondary[muscleID] = true
		}
	}

	// Compare primary muscles
	if len(existingPrimary) != len(importedPrimary) {
		return true
	}
	for muscleID := range existingPrimary {
		if !importedPrimary[muscleID] {
			return true
		}
	}

	// Compare secondary muscles
	if len(existingSecondary) != len(importedSecondary) {
		return true
	}
	for muscleID := range existingSecondary {
		if !importedSecondary[muscleID] {
			return true
		}
	}

	return false
}

// updateMuscleAssociations replaces muscle associations for an exercise
func (db *Database) updateMuscleAssociations(exerciseID uint, importedExercise importedExercise, lookupMap map[string]uint) error {
	exercise := Exercise{ID: exerciseID}

	// Clear existing associations
	if err := db.Model(&exercise).Association("PrimaryMuscles").Clear(); err != nil {
		return fmt.Errorf("failed to clear primary muscles: %w", err)
	}
	if err := db.Model(&exercise).Association("SecondaryMuscles").Clear(); err != nil {
		return fmt.Errorf("failed to clear secondary muscles: %w", err)
	}

	// Add primary muscles
	var primaryMuscles []Muscle
	for _, muscleName := range importedExercise.PrimaryMuscles {
		if muscleID, ok := lookupMap[muscleName]; ok {
			primaryMuscles = append(primaryMuscles, Muscle{ID: muscleID})
		}
	}
	if len(primaryMuscles) > 0 {
		if err := db.Model(&exercise).Association("PrimaryMuscles").Append(&primaryMuscles); err != nil {
			return fmt.Errorf("failed to add primary muscles: %w", err)
		}
	}

	// Add secondary muscles
	var secondaryMuscles []Muscle
	for _, muscleName := range importedExercise.SecondaryMuscles {
		if muscleID, ok := lookupMap[muscleName]; ok {
			secondaryMuscles = append(secondaryMuscles, Muscle{ID: muscleID})
		}
	}
	if len(secondaryMuscles) > 0 {
		if err := db.Model(&exercise).Association("SecondaryMuscles").Append(&secondaryMuscles); err != nil {
			return fmt.Errorf("failed to add secondary muscles: %w", err)
		}
	}

	return nil
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

var (
	idReplacer = strings.NewReplacer(
		" ", "_",
		"/", "_",
		",", "",
		"(", "",
		")", "",
		"-", "-",
		"'", "",
	)
	lastUpdate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
)

func (e Exercise) StringID() string {
	return idReplacer.Replace(e.Name)
}

func (e Exercise) Images() (images []string) {
	id := e.StringID()

	for i := range imageAmount {
		images = append(images, fmt.Sprintf(imageFormat, id, i))
	}
	return
}

func (db *Database) UpdateExercises() (err error) {
	// Load exercises
	exercises, err := downloadExercises()
	if err != nil {
		log.Fatalf("Failed to load exercises: %v", err)
	}

	log.Printf("Successfully loaded %d exercises from JSON", len(exercises))

	// Create/update lookup entities and get maps
	lookupMaps, err := db.upsertLookupEntities(exercises)
	if err != nil {
		return errors.New("Failed to upsert lookup entities: " + err.Error())
	}

	var successCount, createCount, updateCount int

	// Import/update exercises
	for i, exercise := range exercises {
		didSave, isUpdate, err := db.upsertExercise(exercise, lookupMaps)
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
