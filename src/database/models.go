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

const (
	dbDir  = "data"
	dbName = "fitness.db"
)

type Equipment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Exercises []Exercise `gorm:"many2many:exercise_equipment" json:"exercises,omitempty"`
}

type MuscleGroup struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Exercises []Exercise `gorm:"many2many:exercise_muscle_groups" json:"exercises,omitempty"`
}

type ExerciseMuscleGroup struct {
	ID            uint `gorm:"primaryKey" json:"id"`
	ExerciseID    uint `gorm:"uniqueIndex:idx_exercise_muscle_group" json:"exercise_id"`
	MuscleGroupID uint `gorm:"uniqueIndex:idx_exercise_muscle_group" json:"muscle_group_id"`

	Exercise    Exercise    `json:"exercise"`
	MuscleGroup MuscleGroup `json:"muscle_group"`
}

type Exercise struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	//UserID      uint           `gorm:"index" json:"user_id"`
	//IsPublic  bool           `gorm:"default:false" json:"is_public"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	//User         User          `json:"-"`
	Equipment    []Equipment   `gorm:"many2many:exercise_equipment" json:"equipment"`
	MuscleGroups []MuscleGroup `gorm:"many2many:exercise_muscle_groups" json:"muscle_groups"`
	Sets         []Set         `json:"sets,omitempty"`
}

type Set struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ExerciseID uint           `gorm:"index" json:"exercise_id"`
	Reps       int            `json:"reps"`
	Weight     float64        `json:"weight"`
	Duration   int            `json:"duration"` // In seconds, for timed exercises
	OrderIndex int            `gorm:"not null" json:"order_index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Exercise Exercise `json:"-"`
}

// SuperSet to handle two exercises with single rest time
type SuperSet struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Name                string         `gorm:"size:100" json:"name"`
	PrimaryExerciseID   uint           `gorm:"index" json:"primary_exercise_id"`
	SecondaryExerciseID uint           `gorm:"index" json:"secondary_exercise_id"`
	RestTime            int            `gorm:"default:0" json:"rest_time"` // In seconds
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	PrimaryExercise   Exercise `json:"primary_exercise"`
	SecondaryExercise Exercise `json:"secondary_exercise"`
}

// RoutineItem represents either an Exercise or a SuperSet in a Routine
type RoutineItem struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	RoutineID  uint           `gorm:"index" json:"routine_id"`
	ExerciseID *uint          `gorm:"index" json:"exercise_id"`
	SuperSetID *uint          `gorm:"index" json:"super_set_id"`
	RestTime   int            `gorm:"default:0" json:"rest_time"` // In seconds
	OrderIndex int            `gorm:"not null" json:"order_index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Routine  Routine   `json:"-"`
	SuperSet *SuperSet `json:"super_set,omitempty"`
	Exercise *Exercise `json:"exercise,omitempty"`
}

type Routine struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	//UserID      uint           `gorm:"index" json:"user_id"`
	//IsPublic  bool           `gorm:"default:false" json:"is_public"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	//User        User           `json:"-"`
	RoutineItems []RoutineItem `json:"routine_items,omitempty"`
}

/*
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email     string         `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	Name      string         `gorm:"size:100" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Exercises      []Exercise      `json:"exercises,omitempty"`
	Routines       []Routine       `json:"routines,omitempty"`
	RecordRoutines []RecordRoutine `json:"record_routines,omitempty"`
}
*/

type RecordRoutine struct {
	ID uint `gorm:"primaryKey" json:"id"`
	//UserID    uint           `gorm:"index" json:"user_id"`
	RoutineID uint           `gorm:"index" json:"routine_id"`
	StartedAt time.Time      `gorm:"not null" json:"started_at"`
	EndedAt   *time.Time     `json:"ended_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	//User            User             `json:"-"`
	Routine            Routine             `json:"routine"`
	RecordRoutineItems []RecordRoutineItem `json:"record_routine_items,omitempty"`
}

// RecordRoutineItem represents either a RecordExercise or a RecordSuperSet in a completed routine
type RecordRoutineItem struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	RecordRoutineID  uint           `gorm:"index" json:"record_routine_id"`
	RecordExerciseID *uint          `gorm:"index" json:"record_exercise_id"`
	RecordSuperSetID *uint          `gorm:"index" json:"record_super_set_id"`
	ActualRestTime   int            `json:"actual_rest_time"` // In seconds
	OrderIndex       int            `gorm:"not null" json:"order_index"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	RecordRoutine  RecordRoutine   `json:"-"`
	RecordSuperSet *RecordSuperSet `json:"record_super_set,omitempty"`
	RecordExercise *RecordExercise `json:"record_exercise,omitempty"`
}

// RecordSuperSet records a completed superset
type RecordSuperSet struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RecordRoutineID uint           `gorm:"index" json:"record_routine_id"`
	SuperSetID      uint           `gorm:"index" json:"super_set_id"`
	StartedAt       time.Time      `gorm:"not null" json:"started_at"`
	EndedAt         time.Time      `gorm:"not null" json:"ended_at"`
	ActualRestTime  int            `json:"actual_rest_time"` // In seconds
	OrderIndex      int            `gorm:"not null" json:"order_index"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	RecordRoutine RecordRoutine `json:"-"`
	SuperSet      SuperSet      `json:"super_set"`
}

type RecordExercise struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	RecordRoutineID  uint            `gorm:"index" json:"record_routine_id"`
	RecordRoutine    RecordRoutine   `json:"-"`
	ExerciseID       uint            `gorm:"index" json:"exercise_id"`
	Exercise         Exercise        `json:"exercise"`
	StartedAt        time.Time       `gorm:"not null" json:"started_at"`
	EndedAt          time.Time       `gorm:"not null" json:"ended_at"`
	ActualRestTime   int             `json:"actual_rest_time"` // In seconds
	RecordSets       []RecordSet     `json:"record_sets,omitempty"`
	OrderIndex       int             `gorm:"not null" json:"order_index"`
	RecordSuperSetID *uint           `gorm:"index" json:"record_super_set_id"`
	RecordSuperSet   *RecordSuperSet `json:"-"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
}

type RecordSet struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	RecordExerciseID uint           `gorm:"index" json:"record_exercise_id"`
	SetID            uint           `gorm:"index" json:"set_id"`
	ActualReps       int            `json:"actual_reps"`
	ActualWeight     float64        `json:"actual_weight"`
	ActualDuration   int            `json:"actual_duration"` // In seconds
	CompletedAt      time.Time      `gorm:"not null" json:"completed_at"`
	OrderIndex       int            `gorm:"not null" json:"order_index"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	RecordExercise RecordExercise `json:"-"`
	Set            Set            `json:"set"`
}

type Localization struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	LanguageID uint           `gorm:"not null;uniqueIndex:idx_lang_keyword" json:"language_id"`
	Keyword    string         `gorm:"size:255;not null;uniqueIndex:idx_lang_keyword" json:"keyword"`
	Text       string         `gorm:"size:1000;not null" json:"text"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Language struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Code      string         `gorm:"size:8;not null;uniqueIndex" json:"code"`
	Flag      string         `gorm:"size:50" json:"flag"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// InitializeDB creates and initializes the SQLite database with all models
func InitializeDB() (db *gorm.DB, err error) {
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
	db, err = gorm.Open(dialector, config)
	if err != nil {
		return
	}

	// Get the underlying SQL database to set connection parameters
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate the models
	err = db.AutoMigrate(
		Equipment{},
		MuscleGroup{},
		Exercise{},
		ExerciseMuscleGroup{},
		Set{},
		SuperSet{},
		RoutineItem{},
		Routine{},
		//User{},
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

	// Ensure initial data is present
	err = CheckInitialData(db)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
