package database

import (
	"log"

	"gorm.io/gorm"
)

// CheckInitialData ensures that all necessary initial data is in the database
func CheckInitialData(db *gorm.DB) (err error) {
	err = ensureEquipmentData(db)
	if err != nil {
		return
	}

	err = ensureMuscleGroupData(db)
	if err != nil {
		return
	}

	err = ensureExerciseData(db)
	if err != nil {
		return
	}

	log.Println("Initial data verification complete")
	return
}

// ensureEquipmentData checks if equipment data exists and adds it if not
func ensureEquipmentData(db *gorm.DB) error {
	equipmentList := []string{
		"None",
		"Barbell",
		"Dumbbell",
		"Kettlebell",
		"Machine",
		"Plate",
		"ResistanceBand",
		"Suspension",
		"Other",
	}

	// Check if equipment data already exists
	var count int64
	if err := db.Model(&Equipment{}).Count(&count).Error; err != nil {
		return err
	}

	// If no equipment data, insert the initial data
	if count == 0 {
		log.Println("Adding initial equipment data")
		for _, name := range equipmentList {
			equipment := Equipment{
				Name: name,
			}
			if err := db.Create(&equipment).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// ensureMuscleGroupData checks if muscle group data exists and adds it if not
func ensureMuscleGroupData(db *gorm.DB) error {
	muscleGroupList := []string{
		"Abdominals",
		"Abductors",
		"Adductors",
		"Biceps",
		"LowerBack",
		"UpperBack",
		"Cardio",
		"Chest",
		"Calves",
		"Forearms",
		"Glutes",
		"Hamstrings",
		"Lats",
		"Quadriceps",
		"Shoulders",
		"Triceps",
		"Traps",
		"Neck",
		"FullBody",
		"Other",
	}

	// Check if muscle group data already exists
	var count int64
	if err := db.Model(&MuscleGroup{}).Count(&count).Error; err != nil {
		return err
	}

	// If no muscle group data, insert the initial data
	if count == 0 {
		log.Println("Adding initial muscle group data")
		for _, name := range muscleGroupList {
			muscleGroup := MuscleGroup{
				Name: name,
			}
			if err := db.Create(&muscleGroup).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// ensureExerciseData checks if exercise data exists and adds it if not
func ensureExerciseData(db *gorm.DB) error {
	exerciseList := []Exercise{
		{Name: "BenchPress", MuscleGroups: []MuscleGroup{{ID: 8}, {ID: 15}, {ID: 16}}, Equipment: []Equipment{{ID: 2}}},
		{Name: "Squat", MuscleGroups: []MuscleGroup{{ID: 14}, {ID: 12}, {ID: 11}, {ID: 5}}, Equipment: []Equipment{{ID: 2}}},
		{Name: "Deadlift", MuscleGroups: []MuscleGroup{{ID: 5}, {ID: 6}, {ID: 12}, {ID: 11}, {ID: 14}}, Equipment: []Equipment{{ID: 2}}},
		{Name: "OverheadPress", MuscleGroups: []MuscleGroup{{ID: 15}, {ID: 16}}, Equipment: []Equipment{{ID: 2}}},
		{Name: "PullUp", MuscleGroups: []MuscleGroup{{ID: 13}, {ID: 6}, {ID: 4}}, Equipment: []Equipment{{ID: 1}}},
		{Name: "PushUp", MuscleGroups: []MuscleGroup{{ID: 8}, {ID: 15}, {ID: 16}, {ID: 1}}, Equipment: []Equipment{{ID: 1}}},
		{Name: "Lunges", MuscleGroups: []MuscleGroup{{ID: 14}, {ID: 12}, {ID: 11}}, Equipment: []Equipment{{ID: 1}}},
		{Name: "Plank", MuscleGroups: []MuscleGroup{{ID: 1}, {ID: 5}}, Equipment: []Equipment{{ID: 1}}},
		{Name: "BicepCurl", MuscleGroups: []MuscleGroup{{ID: 4}, {ID: 10}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "TricepDip", MuscleGroups: []MuscleGroup{{ID: 16}, {ID: 8}, {ID: 15}}, Equipment: []Equipment{{ID: 1}}},
		{Name: "LegPress", MuscleGroups: []MuscleGroup{{ID: 14}, {ID: 12}, {ID: 11}}, Equipment: []Equipment{{ID: 5}}},
		{Name: "LatPulldown", MuscleGroups: []MuscleGroup{{ID: 13}, {ID: 6}, {ID: 4}}, Equipment: []Equipment{{ID: 5}}},
		{Name: "LegExtension", MuscleGroups: []MuscleGroup{{ID: 14}}, Equipment: []Equipment{{ID: 5}}},
		{Name: "LegCurl", MuscleGroups: []MuscleGroup{{ID: 12}}, Equipment: []Equipment{{ID: 5}}},
		{Name: "ShoulderPress", MuscleGroups: []MuscleGroup{{ID: 15}, {ID: 16}}, Equipment: []Equipment{{ID: 2}}},
		{Name: "ChestFly", MuscleGroups: []MuscleGroup{{ID: 8}}, Equipment: []Equipment{{ID: 5}}},
		{Name: "CableRow", MuscleGroups: []MuscleGroup{{ID: 6}, {ID: 13}, {ID: 4}}, Equipment: []Equipment{{ID: 5}}},
		{Name: "SeatedRow", MuscleGroups: []MuscleGroup{{ID: 6}, {ID: 13}, {ID: 4}}, Equipment: []Equipment{{ID: 5}}},
		{Name: "DumbbellFly", MuscleGroups: []MuscleGroup{{ID: 8}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellRow", MuscleGroups: []MuscleGroup{{ID: 6}, {ID: 13}, {ID: 4}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellShoulderPress", MuscleGroups: []MuscleGroup{{ID: 15}, {ID: 16}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellLateralRaise", MuscleGroups: []MuscleGroup{{ID: 15}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellFrontRaise", MuscleGroups: []MuscleGroup{{ID: 15}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellShrug", MuscleGroups: []MuscleGroup{{ID: 17}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellTricepExtension", MuscleGroups: []MuscleGroup{{ID: 16}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellBicepCurl", MuscleGroups: []MuscleGroup{{ID: 4}, {ID: 10}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellLunge", MuscleGroups: []MuscleGroup{{ID: 14}, {ID: 12}, {ID: 11}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellSquat", MuscleGroups: []MuscleGroup{{ID: 14}, {ID: 12}, {ID: 11}, {ID: 5}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellDeadlift", MuscleGroups: []MuscleGroup{{ID: 5}, {ID: 6}, {ID: 12}, {ID: 11}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellChestPress", MuscleGroups: []MuscleGroup{{ID: 8}, {ID: 15}, {ID: 16}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellChestFly", MuscleGroups: []MuscleGroup{{ID: 8}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellBentOverRow", MuscleGroups: []MuscleGroup{{ID: 6}, {ID: 13}, {ID: 4}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellLateralRaise", MuscleGroups: []MuscleGroup{{ID: 15}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellFrontRaise", MuscleGroups: []MuscleGroup{{ID: 15}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellShrug", MuscleGroups: []MuscleGroup{{ID: 17}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellTricepKickback", MuscleGroups: []MuscleGroup{{ID: 16}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellHammerCurl", MuscleGroups: []MuscleGroup{{ID: 4}, {ID: 10}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellConcentrationCurl", MuscleGroups: []MuscleGroup{{ID: 4}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellSkullCrusher", MuscleGroups: []MuscleGroup{{ID: 16}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellWristCurl", MuscleGroups: []MuscleGroup{{ID: 10}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellWristExtension", MuscleGroups: []MuscleGroup{{ID: 10}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellSideBend", MuscleGroups: []MuscleGroup{{ID: 1}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellRussianTwist", MuscleGroups: []MuscleGroup{{ID: 1}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellPlankRow", MuscleGroups: []MuscleGroup{{ID: 1}, {ID: 6}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellSidePlank", MuscleGroups: []MuscleGroup{{ID: 1}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellMountainClimber", MuscleGroups: []MuscleGroup{{ID: 1}, {ID: 14}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellBicycleCrunch", MuscleGroups: []MuscleGroup{{ID: 1}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellLegRaise", MuscleGroups: []MuscleGroup{{ID: 1}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellReverseCrunch", MuscleGroups: []MuscleGroup{{ID: 1}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellFlutterKick", MuscleGroups: []MuscleGroup{{ID: 1}, {ID: 14}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellSideCrunch", MuscleGroups: []MuscleGroup{{ID: 1}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellToeTouch", MuscleGroups: []MuscleGroup{{ID: 1}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellWoodchopper", MuscleGroups: []MuscleGroup{{ID: 1}, {ID: 5}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellSideLegRaise", MuscleGroups: []MuscleGroup{{ID: 2}, {ID: 11}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellGluteBridge", MuscleGroups: []MuscleGroup{{ID: 11}, {ID: 12}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellCalfRaise", MuscleGroups: []MuscleGroup{{ID: 9}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellStepUp", MuscleGroups: []MuscleGroup{{ID: 14}, {ID: 12}, {ID: 11}}, Equipment: []Equipment{{ID: 3}}},
		{Name: "DumbbellBulgarianSplitSquat", MuscleGroups: []MuscleGroup{{ID: 14}, {ID: 12}, {ID: 11}}, Equipment: []Equipment{{ID: 3}}},
	}

	// Check if exercise data already exists
	var count int64
	if err := db.Model(&Exercise{}).Count(&count).Error; err != nil {
		return err
	}

	// If no exercise data, insert the initial data
	if count == 0 {
		log.Println("Adding initial exercise data")
		for _, exercise := range exerciseList {
			if err := db.Create(&exercise).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
