package database

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50" json:"name"`
	IsFemale  bool           `json:"isFemale"`
	Height    *float64       `json:"height"` // In cm
	Weight    *float64       `json:"weight"` // In kg
	BirthDate *time.Time     `json:"birthDate"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Exercise struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string  `gorm:"not null;uniqueIndex" json:"name"`
	Level        string  `gorm:"size:50;not null" json:"level"`
	Category     string  `gorm:"size:50;not null" json:"category"`
	Force        *string `gorm:"size:50" json:"force"`
	Mechanic     *string `gorm:"size:50" json:"mechanic"`
	Equipment    *string `gorm:"size:50" json:"equipment"`
	Instructions *string `json:"instructions"`

	PrimaryMuscles   []Muscle `gorm:"many2many:exercise_primary_muscles;constraint:OnDelete:CASCADE" json:"primaryMuscles"`
	SecondaryMuscles []Muscle `gorm:"many2many:exercise_secondary_muscles;constraint:OnDelete:CASCADE" json:"secondaryMuscles"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Muscle struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;size:50;not null" json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Routine represents a workout routine blueprint
type Routine struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Items []RoutineItem `json:"items"`
}

// RoutineItem can be either a single exercise or a superset
type RoutineItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoutineID  uint      `gorm:"index;not null" json:"routineId"`
	Type       string    `gorm:"size:20;not null" json:"type"` // "exercise" or "superset"
	RestTime   int       `gorm:"default:0" json:"restTime"`    // In seconds
	OrderIndex int       `gorm:"not null" json:"orderIndex"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	Routine       Routine        `json:"-"`
	ExerciseItems []ExerciseItem `json:"exerciseItems,omitempty"` // For both single exercises and superset items
}

// ExerciseItem represents an exercise within a routine item (could be standalone or part of superset)
type ExerciseItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RoutineItemID uint      `gorm:"index;not null" json:"routineItemId"`
	ExerciseID    uint      `gorm:"index;not null" json:"exerciseId"`
	OrderIndex    int       `gorm:"not null" json:"orderIndex"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`

	RoutineItem RoutineItem `json:"-"`
	Exercise    Exercise    `json:"exercise"`
	Sets        []Set       `json:"sets"`
}

// Set represents a planned set within an exercise
type Set struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ExerciseItemID uint      `gorm:"index;not null" json:"exerciseItemId"`
	Reps           int       `json:"reps"`
	Weight         float64   `json:"weight"`
	Duration       int       `json:"duration"` // In seconds
	OrderIndex     int       `gorm:"not null" json:"orderIndex"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	ExerciseItem ExerciseItem `json:"-"`
}

// ===== RECORD MODELS (for actual workout completion) =====

// RecordRoutine records a completed workout session
type RecordRoutine struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoutineID uint      `gorm:"index;not null" json:"routineId"`
	Duration  *uint     `json:"duration"` // In seconds
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Routine     Routine      `json:"routine"`
	RecordItems []RecordItem `json:"recordItems"`
}

// RecordItem records completion of a routine item (exercise or superset)
type RecordItem struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	RecordRoutineID uint      `gorm:"index;not null" json:"recordRoutineId"`
	RoutineItemID   uint      `gorm:"index;not null" json:"routineItemId"`
	Duration        *uint     `json:"duration"`       // In seconds
	ActualRestTime  *int      `json:"actualRestTime"` // In seconds
	OrderIndex      int       `gorm:"not null" json:"orderIndex"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	RecordRoutine       RecordRoutine        `json:"-"`
	RoutineItem         RoutineItem          `json:"routineItem"`
	RecordExerciseItems []RecordExerciseItem `json:"recordExerciseItems"`
}

// RecordExerciseItem records completion of an exercise within a routine item
type RecordExerciseItem struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	RecordItemID   uint      `gorm:"index;not null" json:"recordItemId"`
	ExerciseItemID uint      `gorm:"index;not null" json:"exerciseItemId"`
	OrderIndex     int       `gorm:"not null" json:"orderIndex"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	RecordItem   RecordItem   `json:"-"`
	ExerciseItem ExerciseItem `json:"exerciseItem"`
	RecordSets   []RecordSet  `json:"recordSets"`
}

// RecordSet records completion of an actual set
type RecordSet struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	RecordExerciseItemID uint      `gorm:"index;not null" json:"recordExerciseItemId"`
	SetID                uint      `gorm:"index;not null" json:"setId"`
	ActualReps           int       `json:"actualReps"`
	ActualWeight         float64   `json:"actualWeight"`
	ActualDuration       int       `json:"actualDuration"` // In seconds
	CompletedAt          time.Time `gorm:"not null" json:"completedAt"`
	OrderIndex           int       `gorm:"not null" json:"orderIndex"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`

	RecordExerciseItem RecordExerciseItem `json:"-"`
	Set                Set                `json:"set"`
}

// InitializeDB creates and initializes the SQLite database with all models
func InitializeDB() (db *Database, err error) {
	// Create the data directory if it doesn't exist
	if _, err = os.Stat(dbDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbDir, 0755)
		if err != nil {
			return
		}
	}

	dbPath := filepath.Join(dbDir, dbName)

	// Set up logger for GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	dialector := sqlite.Open(dbPath + "?_pragma=foreign_keys(1)")
	config := &gorm.Config{Logger: newLogger}

	// Open connection to the database
	conn, err := gorm.Open(dialector, config)
	if err != nil {
		return
	}

	// Get the underlying SQL database to set connection parameters
	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate the models in correct order
	err = conn.AutoMigrate(
		&User{},
		&Muscle{},
		&Exercise{},
		&Routine{},
		&RoutineItem{},
		&ExerciseItem{},
		&Set{},
		&RecordRoutine{},
		&RecordItem{},
		&RecordExerciseItem{},
		&RecordSet{},
	)
	if err != nil {
		return
	}

	db = &Database{conn}

	// Ensure initial data is present
	err = db.CheckInitialData()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Helper methods for creating and querying routines

// CreateRoutineWithData creates a routine with all nested data
func (db *Database) CreateRoutineWithData(routine *Routine) error {
	return db.Create(routine).Error
}

// GetRoutineWithItems retrieves a routine with all its nested data
func (db *Database) GetRoutineWithItems(routineID uint) (*Routine, error) {
	var routine Routine
	err := db.Preload("Items.ExerciseItems.Exercise.PrimaryMuscles").
		Preload("Items.ExerciseItems.Exercise.SecondaryMuscles").
		Preload("Items.ExerciseItems.Sets").
		Order("Items.order_index, Items.ExerciseItems.order_index, Items.ExerciseItems.Sets.order_index").
		First(&routine, routineID).Error

	return &routine, err
}

// GetRecordRoutineWithData retrieves a completed workout with all nested data
func (db *Database) GetRecordRoutineWithData(recordID uint) (*RecordRoutine, error) {
	var record RecordRoutine
	err := db.Preload("Routine").
		Preload("RecordItems.RoutineItem").
		Preload("RecordItems.RecordExerciseItems.ExerciseItem.Exercise.PrimaryMuscles").
		Preload("RecordItems.RecordExerciseItems.ExerciseItem.Exercise.SecondaryMuscles").
		Preload("RecordItems.RecordExerciseItems.RecordSets.Set").
		Order("RecordItems.order_index, RecordItems.RecordExerciseItems.order_index, RecordItems.RecordExerciseItems.RecordSets.order_index").
		First(&record, recordID).Error

	return &record, err
}
