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

func (db *Database) GetExerciseByID(id string) (*Exercise, error) {
	var exercise Exercise
	err := db.First(&exercise, "id = ?", id).Error
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

func (db *Database) GetRoutines() ([]Routine, error) {
	var routines []Routine
	err := db.
		Preload("RoutineItems").
		Preload("RoutineItems.ExerciseItems").
		Preload("RoutineItems.ExerciseItems.Exercise").
		Preload("RoutineItems.ExerciseItems.Sets").
		Find(&routines).Error
	if err != nil {
		return nil, err
	}

	return routines, nil
}

func (db *Database) GetRoutineByID(id uint) (*Routine, error) {
	var routine Routine
	err := db.
		Preload("RoutineItems").
		Preload("RoutineItems.ExerciseItems").
		Preload("RoutineItems.ExerciseItems.Exercise").
		Preload("RoutineItems.ExerciseItems.Sets").
		First(&routine, id).Error
	if err != nil {
		return nil, err
	}

	return &routine, nil
}

func (db *Database) NewRoutine(routine *Routine) error {
	if routine.Name == "" || len(routine.Name) > 100 {
		return fmt.Errorf("invalid routine name")
	}

	if err := db.Create(routine).Error; err != nil {
		return fmt.Errorf("failed to create routine: %w", err)
	}

	return nil
}

func (db *Database) UpdateRoutine(routine *Routine) error {
	if routine.ID == 0 {
		return fmt.Errorf("routine ID is required for update")
	}

	if routine.Name == "" || len(routine.Name) > 100 {
		return fmt.Errorf("invalid routine name")
	}

	if err := db.Save(routine).Error; err != nil {
		return fmt.Errorf("failed to update routine: %w", err)
	}

	return nil
}

func (db *Database) DeleteRoutine(id uint) error {
	if id == 0 {
		return fmt.Errorf("routine ID is required for deletion")
	}

	if err := db.Delete(&Routine{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete routine: %w", err)
	}

	return nil
}

func (db *Database) NewRoutineItem(r *Routine) (*RoutineItem, error) {
	item := &RoutineItem{
		Routine:    *r,
		OrderIndex: len(r.RoutineItems),
	}

	err := db.Create(item).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create new routine item: %w", err)
	}

	return item, nil
}

func (db *Database) GetRoutineItemByID(id uint) (*RoutineItem, error) {
	var item RoutineItem
	err := db.
		Preload("ExerciseItems").
		Preload("ExerciseItems.Exercise").
		First(&item, id).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (db *Database) DeleteRoutineItem(item *RoutineItem) (uint, error) {
	if item.ID == 0 {
		return 0, fmt.Errorf("routine item ID is required for deletion")
	}

	if err := db.Delete(item).Error; err != nil {
		return 0, fmt.Errorf("failed to delete routine item: %w", err)
	}

	return item.RoutineID, nil
}

func (db *Database) NewExerciseItem(routineItem *RoutineItem, exerciseID string) (*ExerciseItem, error) {
	item := &ExerciseItem{
		RoutineItemID: routineItem.ID,
		ExerciseID:    exerciseID,
		OrderIndex:    len(routineItem.ExerciseItems),
	}

	if err := db.Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to create new exercise item: %w", err)
	}

	emptySet := &Set{ExerciseItemID: item.ID}
	if err := db.Create(emptySet).Error; err != nil {
		return nil, fmt.Errorf("failed to create initial set for exercise item: %w", err)
	}
	item.Sets = append(item.Sets, *emptySet)

	return item, nil
}

func (db *Database) DeleteExerciseItem(item *ExerciseItem) (uint, error) {
	if item.ID == 0 {
		return 0, fmt.Errorf("exercise item ID is required for deletion")
	}

	if err := db.Delete(item).Error; err != nil {
		return 0, fmt.Errorf("failed to delete exercise item: %w", err)
	}

	routineItem, err := db.GetRoutineItemByID(item.RoutineItemID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve routine item: %w", err)
	}

	if len(routineItem.ExerciseItems) == 0 {
		if err := db.Delete(routineItem).Error; err != nil {
			return 0, fmt.Errorf("failed to delete routine item after removing last exercise item: %w", err)
		}
	}

	return routineItem.RoutineID, nil
}

func (db *Database) GetExerciseItemByID(id uint) (*ExerciseItem, error) {
	var item ExerciseItem
	err := db.
		Preload("Exercise").
		Preload("Sets").
		Preload("RoutineItem").
		First(&item, id).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (db *Database) AddExerciseToRoutineItem(routineItem *RoutineItem, exercise *Exercise) (*ExerciseItem, error) {
	item := &ExerciseItem{
		RoutineItemID: routineItem.ID,
		ExerciseID:    exercise.ID,
		OrderIndex:    len(routineItem.ExerciseItems),
	}

	if err := db.Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to add exercise to routine item: %w", err)
	}

	emptySet := &Set{ExerciseItemID: item.ID}
	if err := db.Create(emptySet).Error; err != nil {
		return nil, fmt.Errorf("failed to create initial set for exercise item: %w", err)
	}
	item.Sets = append(item.Sets, *emptySet)

	return item, nil
}

func (db *Database) UpdateExerciseItem(item *ExerciseItem) error {
	if item.ID == 0 {
		return fmt.Errorf("exercise item ID is required for update")
	}

	if item.RestTime > 3600 {
		return fmt.Errorf("invalid rest time: %d", item.RestTime)
	}

	if item.ExerciseID == "" {
		return fmt.Errorf("exercise ID is required")
	}

	if err := db.Save(item).Error; err != nil {
		return fmt.Errorf("failed to update exercise item: %w", err)
	}

	return nil
}

func (db *Database) NewSet(item *ExerciseItem) (*Set, error) {
	if item.ID == 0 {
		return nil, fmt.Errorf("exercise item ID is required for new set")
	}

	l := len(item.Sets)

	set := &Set{
		ExerciseItemID: item.ID,
		OrderIndex:     l,
	}

	if l > 0 {
		lastSet := item.Sets[l-1]

		set.Weight = lastSet.Weight
		set.Reps = lastSet.Reps
		set.Duration = lastSet.Duration
	}

	if err := db.Create(set).Error; err != nil {
		return nil, fmt.Errorf("failed to create new set: %w", err)
	}

	return set, nil
}

func (db *Database) GetSetByID(id uint) (*Set, error) {
	var set Set
	err := db.
		Preload("ExerciseItem").
		Preload("ExerciseItem.RoutineItem").
		Preload("ExerciseItem.Exercise").
		First(&set, id).Error
	if err != nil {
		return nil, err
	}

	return &set, nil
}

func (db *Database) UpdateSet(set *Set) error {
	if set.ID == 0 {
		return fmt.Errorf("set ID is required for update")
	}

	// Check if reps is valid
	if set.Reps != nil && (*set.Reps == 0 || *set.Reps > 99) {
		return fmt.Errorf("invalid reps value: %v", *set.Reps)
	}

	// Check if weight is valid
	if set.Weight != nil && (*set.Weight <= 0 || *set.Weight > 300) {
		return fmt.Errorf("invalid weight value: %v", *set.Weight)
	}

	// Check if duration is valid
	if set.Duration != nil && (*set.Duration == 0 || *set.Duration > 7200) {
		return fmt.Errorf("invalid duration value: %v", *set.Duration)
	}

	if err := db.Save(set).Error; err != nil {
		return fmt.Errorf("failed to update set: %w", err)
	}

	return nil
}

func (db *Database) DeleteSet(set *Set) (uint, error) {
	if set.ID == 0 {
		return 0, fmt.Errorf("set ID is required for deletion")
	}

	if err := db.Delete(set).Error; err != nil {
		return 0, fmt.Errorf("failed to delete set: %w", err)
	}

	return set.ExerciseItem.RoutineItem.RoutineID, nil
}
