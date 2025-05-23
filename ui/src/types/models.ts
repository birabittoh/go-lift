// Models for Go-Lift app that match the backend models

// Equipment model
export interface Equipment {
  id?: number;
  name: string;
  description: string;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  exercises?: Exercise[]; // Many-to-many relationship
}

// Muscle Group model
export interface MuscleGroup {
  id?: number;
  name: string;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  exercises?: Exercise[]; // Many-to-many relationship
}

// ExerciseMuscleGroup join table
export interface ExerciseMuscleGroup {
  id?: number;
  exerciseId: number;
  muscleGroupId: number;
  exercise?: Exercise;
  muscleGroup?: MuscleGroup;
}

// Set definition for an exercise
export interface Set {
  id?: number;
  exerciseId: number;
  reps: number;
  weight: number;
  duration: number; // In seconds, for timed exercises
  orderIndex: number;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  exercise?: Exercise; // Excluded in JSON via json:"-" but useful for frontend
}

// Exercise definition
export interface Exercise {
  id?: number;
  name: string;
  description: string;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  equipment: Equipment[];
  muscleGroups: MuscleGroup[];
  sets?: Set[];
}

// SuperSet to handle two exercises with single rest time
export interface SuperSet {
  id?: number;
  name: string;
  primaryExerciseId: number;
  secondaryExerciseId: number;
  restTime: number; // In seconds
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  primaryExercise: Exercise;
  secondaryExercise: Exercise;
}

// RoutineItem represents either an Exercise or a SuperSet in a Routine
export interface RoutineItem {
  id?: number;
  routineId: number;
  exerciseId?: number | null;
  superSetId?: number | null;
  restTime: number; // In seconds
  orderIndex: number;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  superSet?: SuperSet | null;
  exercise?: Exercise | null;
}

// A collection of exercises and/or supersets that make up a workout routine
export interface Routine {
  id?: number;
  name: string;
  description: string;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  routineItems: RoutineItem[];
}

// RecordRoutine represents a completed workout session
export interface RecordRoutine {
  id?: number;
  routineId: number;
  startedAt: string;
  endedAt?: string | null;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  routine: Routine;
  recordRoutineItems: RecordRoutineItem[];
}

// RecordRoutineItem represents either a RecordExercise or a RecordSuperSet in a completed routine
export interface RecordRoutineItem {
  id?: number;
  recordRoutineId: number;
  recordExerciseId?: number | null;
  recordSuperSetId?: number | null;
  actualRestTime: number; // In seconds
  orderIndex: number;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  recordSuperSet?: RecordSuperSet | null;
  recordExercise?: RecordExercise | null;
}

// RecordSuperSet records a completed superset
export interface RecordSuperSet {
  id?: number;
  recordRoutineId: number;
  superSetId: number;
  startedAt: string;
  endedAt: string;
  actualRestTime: number; // In seconds
  orderIndex: number;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  superSet: SuperSet;
}

// RecordExercise tracks a completed exercise
export interface RecordExercise {
  id?: number;
  recordRoutineId: number;
  exerciseId: number;
  startedAt: string;
  endedAt: string;
  actualRestTime: number; // In seconds
  orderIndex: number;
  recordSuperSetId?: number | null;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  exercise: Exercise;
  recordSets: RecordSet[];
}

// RecordSet tracks an individual completed set
export interface RecordSet {
  id?: number;
  recordExerciseId: number;
  setId: number;
  actualReps: number;
  actualWeight: number;
  actualDuration: number; // In seconds
  completedAt: string;
  orderIndex: number;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
  set: Set;
}

// Additional models for localization
export interface Localization {
  id?: number;
  languageId: number;
  keyword: string;
  text: string;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
}

export interface Language {
  id?: number;
  name: string;
  code: string;
  flag: string;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
}

export interface User {
  id?: number;
  name: string;
  isFemale: string | boolean;
  weight: string | number; // in kg
  height: string | number; // in cm
  birthDate: string; // ISO format date string
  createdAt?: string;
  updatedAt?: string;
}

// Stats for the home page
export interface WorkoutStats {
  totalWorkouts: number;
  totalMinutes: number;
  totalExercises: number;
  mostFrequentExercise?: {
    name: string;
    count: number;
  };
  mostFrequentRoutine?: {
    name: string;
    count: number;
  };
  recentWorkouts: RecordRoutine[];
}

// Some simpler interfaces for UI use when full objects are too complex
export interface ExerciseSimple {
  id?: number;
  name: string;
  muscleGroups: string[];
  equipment: string[];
}

export interface SetForUI {
  id?: number;
  reps: number;
  weight: number;
  duration?: number;
  completed: boolean;
}

export interface ExerciseForWorkout {
  id?: number;
  exerciseId: number;
  exercise?: Exercise;
  sets: SetForUI[];
  notes?: string;
}
