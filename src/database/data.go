package database

import (
	"log"
)

var (
	defaultUserList = []User{{Name: "User"}}
	weekDays        = []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
)

// CheckInitialData ensures that all necessary initial data is in the database
func (db *Database) CheckInitialData() (err error) {
	err = db.ensureExerciseData()
	if err != nil {
		return
	}

	err = db.ensureUserData()
	if err != nil {
		return
	}

	log.Println("Initial data verification complete")
	return
}

// ensureExerciseData checks if exercise data exists and adds it if not
func (db *Database) ensureExerciseData() error {
	// Check if exercise data already exists
	var count int64
	if err := db.Model(&Exercise{}).Count(&count).Error; err != nil {
		return err
	}

	// If no exercise data, insert the initial data
	if count == 0 {
		log.Println("Adding initial exercise data")
		db.UpdateExercises()
	}

	return nil
}

// ensureUserData checks if user data exists and adds it if not
func (db *Database) ensureUserData() error {
	// Check if user data already exists
	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return err
	}

	// If no user data, insert the initial data
	if count == 0 {
		log.Println("Adding initial user data")
		for _, user := range defaultUserList {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
