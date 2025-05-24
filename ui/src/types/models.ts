export interface User {
  id?: number;
  name: string;
  isFemale: boolean;
  height?: number; // In cm
  weight?: number; // In kg
  birthDate?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface Muscle {
  id: number;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface Exercise {
  id: number;
  name: string;
  level: string;
  category: string;
  force?: string;
  mechanic?: string;
  equipment?: string;
  instructions?: string;
  primaryMuscles: Muscle[];
  secondaryMuscles: Muscle[];
  createdAt: string;
  updatedAt: string;
}

export interface Set {
  id: number;
  exerciseItemId: number;
  reps: number;
  weight: number;
  duration: number; // In seconds
  orderIndex: number;
  createdAt: string;
  updatedAt: string;
}

export interface ExerciseItem {
  id: number;
  routineItemId: number;
  exerciseId: number;
  orderIndex: number;
  createdAt: string;
  updatedAt: string;
  exercise: Exercise;
  sets: Set[];
}

export interface RoutineItem {
  id: number;
  routineId: number;
  type: string; // "exercise" or "superset"
  restTime: number; // In seconds
  orderIndex: number;
  createdAt: string;
  updatedAt: string;
  exerciseItems: ExerciseItem[];
}

export interface Routine {
  id: number;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
  items: RoutineItem[];
}

export interface RecordSet {
  id: number;
  recordExerciseItemId: number;
  setId: number;
  actualReps: number;
  actualWeight: number;
  actualDuration: number; // In seconds
  completedAt: string;
  orderIndex: number;
  createdAt: string;
  updatedAt: string;
  set: Set;
}

export interface RecordExerciseItem {
  id: number;
  recordItemId: number;
  exerciseItemId: number;
  orderIndex: number;
  createdAt: string;
  updatedAt: string;
  exerciseItem: ExerciseItem;
  recordSets: RecordSet[];
}

export interface RecordItem {
  id: number;
  recordRoutineId: number;
  routineItemId: number;
  duration?: number; // In seconds
  actualRestTime?: number; // In seconds
  orderIndex: number;
  createdAt: string;
  updatedAt: string;
  routineItem: RoutineItem;
  recordExerciseItems: RecordExerciseItem[];
}

export interface RecordRoutine {
  id: number;
  routineId: number;
  duration?: number; // In seconds
  createdAt: string;
  updatedAt: string;
  routine: Routine;
  recordItems: RecordItem[];
}

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
