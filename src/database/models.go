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

// User model - kept as is since it's not directly related to routines
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

// Measurement models - kept as is
type HeightMeasurement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Height    float64   `json:"height"` // In cm
	CreatedAt time.Time `json:"createdAt"`
}

type WeightMeasurement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Weight    float64   `json:"weight"` // In kg
	CreatedAt time.Time `json:"createdAt"`
}

// Exercise model
type Exercise struct {
	ID                 string  `gorm:"primaryKey" json:"id"`
	Name               string  `gorm:"not null;uniqueIndex" json:"name"`
	Level              string  `gorm:"size:50;not null" json:"level"`
	Category           string  `gorm:"size:50;not null" json:"category"`
	Force              *string `gorm:"size:50" json:"force"`
	Mechanic           *string `gorm:"size:50" json:"mechanic"`
	Equipment          *string `gorm:"size:50" json:"equipment"`
	InstructionsString *string `json:"-"`

	PrimaryMuscles   *string `json:"primaryMuscles"`
	SecondaryMuscles *string `json:"secondaryMuscles"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Non-persisted fields
	Images       []string `gorm:"-" json:"images,omitempty"`
	Instructions []string `gorm:"-" json:"instructions,omitempty"`
}

// Routine represents a workout routine blueprint
type Routine struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	RoutineItems []RoutineItem `gorm:"foreignKey:RoutineID" json:"routineItems"`
}

// RoutineItem represents a group of exercises (can be a single exercise or superset)
type RoutineItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoutineID  uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"routineId"`
	OrderIndex int       `gorm:"not null;default:0" json:"orderIndex"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	Routine       Routine        `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	ExerciseItems []ExerciseItem `gorm:"foreignKey:RoutineItemID" json:"exerciseItems"`
}

// ExerciseItem represents an exercise within a routine item
type ExerciseItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RoutineItemID uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"routineItemId"`
	ExerciseID    string    `gorm:"not null;constraint:OnDelete:CASCADE" json:"exerciseId"`
	RestTime      uint      `gorm:"not null;default:0" json:"restTime"` // In seconds
	OrderIndex    int       `gorm:"not null;default:0" json:"orderIndex"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`

	RoutineItem RoutineItem `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Exercise    Exercise    `gorm:"constraint:OnDelete:CASCADE" json:"exercise"`
	Sets        []Set       `gorm:"foreignKey:ExerciseItemID" json:"sets"`
}

// Set represents a planned set within an exercise
type Set struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ExerciseItemID uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"exerciseItemId"`
	Reps           uint      `gorm:"default:0" json:"reps"`
	Weight         float64   `gorm:"default:0" json:"weight"`
	Duration       uint      `gorm:"default:0" json:"duration"` // In seconds
	OrderIndex     int       `gorm:"not null;default:0" json:"orderIndex"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	ExerciseItem ExerciseItem `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

// ===== RECORD MODELS (for actual workout completion) =====

// RecordRoutine records a completed workout session
type RecordRoutine struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoutineID uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"routineId"`
	Duration  *uint     `json:"duration"` // In seconds, total workout time
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Routine            Routine             `gorm:"constraint:OnDelete:CASCADE" json:"routine"`
	RecordRoutineItems []RecordRoutineItem `gorm:"foreignKey:RecordRoutineID" json:"recordRoutineItems"`
}

// RecordRoutineItem records completion of a routine item (group of exercises)
type RecordRoutineItem struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	RecordRoutineID uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"recordRoutineId"`
	RoutineItemID   uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"routineItemId"`
	OrderIndex      int       `gorm:"not null;default:0" json:"orderIndex"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	RecordRoutine       RecordRoutine        `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	RoutineItem         RoutineItem          `gorm:"constraint:OnDelete:CASCADE" json:"routineItem"`
	RecordExerciseItems []RecordExerciseItem `gorm:"foreignKey:RecordRoutineItemID" json:"recordExerciseItems"`
}

// RecordExerciseItem records completion of an exercise within a routine item
type RecordExerciseItem struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	RecordRoutineItemID uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"recordRoutineItemId"`
	ExerciseItemID      uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"exerciseItemId"`
	ActualRestTime      *int      `json:"actualRestTime"` // In seconds, actual rest taken after this exercise
	OrderIndex          int       `gorm:"not null;default:0" json:"orderIndex"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`

	RecordRoutineItem RecordRoutineItem `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	ExerciseItem      ExerciseItem      `gorm:"constraint:OnDelete:CASCADE" json:"exerciseItem"`
	RecordSets        []RecordSet       `gorm:"foreignKey:RecordExerciseItemID" json:"recordSets"`
}

// RecordSet records completion of an actual set
type RecordSet struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	RecordExerciseItemID uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"recordExerciseItemId"`
	SetID                uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"setId"`
	ActualReps           int       `gorm:"default:0" json:"actualReps"`
	ActualWeight         float64   `gorm:"default:0" json:"actualWeight"`
	ActualDuration       int       `gorm:"default:0" json:"actualDuration"` // In seconds
	CompletedAt          time.Time `gorm:"not null" json:"completedAt"`
	OrderIndex           int       `gorm:"not null;default:0" json:"orderIndex"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`

	RecordExerciseItem RecordExerciseItem `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Set                Set                `gorm:"constraint:OnDelete:CASCADE" json:"set"`
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
		&HeightMeasurement{},
		&WeightMeasurement{},
		&Exercise{},
		&Routine{},
		&RoutineItem{},
		&ExerciseItem{},
		&Set{},
		&RecordRoutine{},
		&RecordRoutineItem{},
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
