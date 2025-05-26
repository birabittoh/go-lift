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
	if user.BirthDate != nil && (user.BirthDate.After(time.Now()) || user.BirthDate.Year() < 1900) {
		return fmt.Errorf("invalid birth date: %v", user.BirthDate)
	}

	if user.Name == "" || len(user.Name) > 50 {
		return fmt.Errorf("invalid name")
	}

	var nw *WeightMeasurement
	var nh *HeightMeasurement

	if user.Weight != nil {
		if *user.Weight < 0 || *user.Weight > 200 {
			return fmt.Errorf("invalid weight: %f", *user.Weight)
		}
		nw = &WeightMeasurement{Weight: *user.Weight}
	}

	if user.Height != nil {
		if *user.Height < 0 || *user.Height > 250 {
			return fmt.Errorf("invalid height: %f", *user.Height)
		}
		nh = &HeightMeasurement{Height: *user.Height}
	}

	if nh != nil {
		var lastHeight HeightMeasurement
		db.Last(&lastHeight)
		if lastHeight.Height != nh.Height {
			if err := db.Create(nh).Error; err != nil {
				return fmt.Errorf("failed to save height measurement: %w", err)
			}
		}
	}

	if nw != nil {
		var lastWeight WeightMeasurement
		db.Last(&lastWeight)
		if lastWeight.Weight != nw.Weight {
			if err := db.Create(nw).Error; err != nil {
				return fmt.Errorf("failed to save weight measurement: %w", err)
			}
		}
	}

	return db.Save(user).Error
}
