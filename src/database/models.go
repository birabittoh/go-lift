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
	IsFemale  bool           `gorm:"default:false" json:"isFemale"`
	Height    float64        `json:"height"` // In cm
	Weight    float64        `json:"weight"` // In kg
	BirthDate time.Time      `json:"birthDate"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Exercise struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	Name         string  `gorm:"not null;uniqueIndex"`
	Level        string  `gorm:"size:50;not null"`
	Category     string  `gorm:"size:50;not null"`
	Force        *string `gorm:"size:50"`
	Mechanic     *string `gorm:"size:50"`
	Equipment    *string `gorm:"size:50"`
	Instructions *string

	PrimaryMuscles   []Muscle `gorm:"many2many:exercise_primary_muscles;constraint:OnDelete:CASCADE"`
	SecondaryMuscles []Muscle `gorm:"many2many:exercise_secondary_muscles;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Muscle struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;size:50;not null"`
}

type Set struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ExerciseID uint           `gorm:"index" json:"exerciseId"`
	Reps       int            `json:"reps"`
	Weight     float64        `json:"weight"`
	Duration   int            `json:"duration"` // In seconds, for timed exercises
	OrderIndex int            `gorm:"not null" json:"orderIndex"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	Exercise Exercise `json:"-"`
}

// SuperSet to handle two exercises with single rest time
type SuperSet struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Name                string         `gorm:"size:100" json:"name"`
	PrimaryExerciseID   uint           `gorm:"index" json:"primaryExerciseId"`
	SecondaryExerciseID uint           `gorm:"index" json:"secondaryExerciseId"`
	RestTime            int            `gorm:"default:0" json:"restTime"` // In seconds
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	PrimaryExercise   Exercise `json:"primaryExercise"`
	SecondaryExercise Exercise `json:"secondaryExercise"`
}

// RoutineItem represents either an Exercise or a SuperSet in a Routine
type RoutineItem struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	RoutineID  uint           `gorm:"index" json:"routineId"`
	ExerciseID *uint          `gorm:"index" json:"exerciseId"`
	SuperSetID *uint          `gorm:"index" json:"superSetId"`
	RestTime   int            `gorm:"default:0" json:"restTime"` // In seconds
	OrderIndex int            `gorm:"not null" json:"orderIndex"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	Routine  Routine   `json:"-"`
	SuperSet *SuperSet `json:"superSet,omitempty"`
	Exercise *Exercise `json:"exercise,omitempty"`
}

type Routine struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	//UserID      uint           `gorm:"index" json:"userId"`
	//IsPublic  bool           `gorm:"default:false" json:"isPublic"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	//User        User           `json:"-"`
	RoutineItems []RoutineItem `json:"routineItems,omitempty"`
}

/*
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email     string         `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	Name      string         `gorm:"size:100" json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	Exercises      []Exercise      `json:"exercises,omitempty"`
	Routines       []Routine       `json:"routines,omitempty"`
	RecordRoutines []RecordRoutine `json:"recordRoutines,omitempty"`
}
*/

type RecordRoutine struct {
	ID uint `gorm:"primaryKey" json:"id"`
	//UserID    uint           `gorm:"index" json:"userId"`
	RoutineID uint           `gorm:"index" json:"routineId"`
	StartedAt time.Time      `gorm:"not null" json:"startedAt"`
	EndedAt   *time.Time     `json:"endedAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	//User            User             `json:"-"`
	Routine            Routine             `json:"routine"`
	RecordRoutineItems []RecordRoutineItem `json:"recordRoutineItems,omitempty"`
}

// RecordRoutineItem represents either a RecordExercise or a RecordSuperSet in a completed routine
type RecordRoutineItem struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	RecordRoutineID  uint           `gorm:"index" json:"recordRoutineId"`
	RecordExerciseID *uint          `gorm:"index" json:"recordExerciseId"`
	RecordSuperSetID *uint          `gorm:"index" json:"recordSuperSetId"`
	ActualRestTime   int            `json:"actualRestTime"` // In seconds
	OrderIndex       int            `gorm:"not null" json:"orderIndex"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	RecordRoutine  RecordRoutine   `json:"-"`
	RecordSuperSet *RecordSuperSet `json:"recordSuperSet,omitempty"`
	RecordExercise *RecordExercise `json:"recordExercise,omitempty"`
}

// RecordSuperSet records a completed superset
type RecordSuperSet struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RecordRoutineID uint           `gorm:"index" json:"recordRoutineId"`
	SuperSetID      uint           `gorm:"index" json:"superSetId"`
	StartedAt       time.Time      `gorm:"not null" json:"startedAt"`
	EndedAt         time.Time      `gorm:"not null" json:"endedAt"`
	ActualRestTime  int            `json:"actualRestTime"` // In seconds
	OrderIndex      int            `gorm:"not null" json:"orderIndex"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	RecordRoutine RecordRoutine `json:"-"`
	SuperSet      SuperSet      `json:"superSet"`
}

type RecordExercise struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	RecordRoutineID  uint            `gorm:"index" json:"recordRoutineId"`
	RecordRoutine    RecordRoutine   `json:"-"`
	ExerciseID       uint            `gorm:"index" json:"exerciseId"`
	Exercise         Exercise        `json:"exercise"`
	StartedAt        time.Time       `gorm:"not null" json:"startedAt"`
	EndedAt          time.Time       `gorm:"not null" json:"endedAt"`
	ActualRestTime   int             `json:"actualRestTime"` // In seconds
	RecordSets       []RecordSet     `json:"recordSets,omitempty"`
	OrderIndex       int             `gorm:"not null" json:"orderIndex"`
	RecordSuperSetID *uint           `gorm:"index" json:"recordSuperSetId"`
	RecordSuperSet   *RecordSuperSet `json:"-"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"deletedAt"`
}

type RecordSet struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	RecordExerciseID uint           `gorm:"index" json:"recordExerciseId"`
	SetID            uint           `gorm:"index" json:"setId"`
	ActualReps       int            `json:"actualReps"`
	ActualWeight     float64        `json:"actualWeight"`
	ActualDuration   int            `json:"actualDuration"` // In seconds
	CompletedAt      time.Time      `gorm:"not null" json:"completedAt"`
	OrderIndex       int            `gorm:"not null" json:"orderIndex"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	RecordExercise RecordExercise `json:"-"`
	Set            Set            `json:"set"`
}

type Localization struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	LanguageID uint           `gorm:"not null;uniqueIndex:idxLangKeyword" json:"languageId"`
	Keyword    string         `gorm:"size:255;not null;uniqueIndex:idxLangKeyword" json:"keyword"`
	Text       string         `gorm:"size:1000;not null" json:"text"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Language struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Code      string         `gorm:"size:8;not null;uniqueIndex" json:"code"`
	Flag      string         `gorm:"size:50" json:"flag"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
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

	dialector := sqlite.Open(dbPath + "?Pragma=foreignKeys(1)")
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

	// Auto migrate the models
	err = conn.AutoMigrate(
		Muscle{},
		Exercise{},
		Set{},
		SuperSet{},
		RoutineItem{},
		Routine{},
		User{},
		RecordRoutine{},
		RecordExercise{},
		RecordSuperSet{},
		RecordSet{},
		Localization{},
		Language{},
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

	log.Println("Database initialized successfully")
	return db, nil
}
