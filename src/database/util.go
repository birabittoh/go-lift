package database

import (
	"fmt"
	"time"
)

func (db *Database) GetExercises() ([]Exercise, error) {
	var exercises []Exercise
	err := db.Find(&exercises).Error
	if err != nil {
		return nil, err
	}

	for i := range exercises {
		exercises[i].Fill()
	}

	return exercises, nil
}

func (db *Database) GetExerciseByID(id uint) (*Exercise, error) {
	var exercise Exercise
	err := db.First(&exercise, id).Error
	if err != nil {
		return nil, err
	}

	exercise.Fill()

	return &exercise, nil
}

func (db *Database) GetUserByID(id uint) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) UpdateUser(user *User) error {
	// validate user.weight, user.height, user.birthdate
	if user.Weight != nil && (*user.Weight < 0 || *user.Weight > 200) {
		return fmt.Errorf("invalid weight: %f", *user.Weight)
	}

	if user.Height != nil && (*user.Height < 0 || *user.Height > 250) {
		return fmt.Errorf("invalid height: %f", *user.Height)
	}

	if user.BirthDate != nil && (user.BirthDate.After(time.Now()) || user.BirthDate.Year() < 1900) {
		return fmt.Errorf("invalid birth date: %v", user.BirthDate)
	}

	if user.Name == "" || len(user.Name) > 50 {
		return fmt.Errorf("invalid name")
	}

	return db.Save(user).Error
}
