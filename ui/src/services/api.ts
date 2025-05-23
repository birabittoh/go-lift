import type { 
  Exercise, 
  Routine, 
  RecordRoutine, 
  User,
  Equipment,
  MuscleGroup,
  Set,
  SuperSet,
  RecordExercise, 
  WorkoutStats
} from '../types/models';

// Base API URL - should be configurable via environment variables in a real app
const API_BASE_URL = '/api';

// Generic fetch with error handling
async function fetchApi<T>(
  endpoint: string, 
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || `API request failed with status ${response.status}`);
  }

  return response.json();
}

// Equipment API services
export const EquipmentService = {
  getAll: () => fetchApi<Equipment[]>('/equipment'),
  
  getById: (id: number) => fetchApi<Equipment>(`/equipment/${id}`),
  
  create: (equipment: Equipment) => fetchApi<Equipment>('/equipment', {
    method: 'POST',
    body: JSON.stringify(equipment),
  }),
  
  update: (id: number, equipment: Equipment) => fetchApi<Equipment>(`/equipment/${id}`, {
    method: 'PUT',
    body: JSON.stringify(equipment),
  }),
  
  delete: (id: number) => fetchApi<void>(`/equipment/${id}`, {
    method: 'DELETE',
  }),
};

// MuscleGroup API services
export const MuscleGroupService = {
  getAll: () => fetchApi<MuscleGroup[]>('/musclegroups'),
  
  getById: (id: number) => fetchApi<MuscleGroup>(`/musclegroups/${id}`),
  
  create: (muscleGroup: MuscleGroup) => fetchApi<MuscleGroup>('/musclegroups', {
    method: 'POST',
    body: JSON.stringify(muscleGroup),
  }),
  
  update: (id: number, muscleGroup: MuscleGroup) => fetchApi<MuscleGroup>(`/musclegroups/${id}`, {
    method: 'PUT',
    body: JSON.stringify(muscleGroup),
  }),
  
  delete: (id: number) => fetchApi<void>(`/musclegroups/${id}`, {
    method: 'DELETE',
  }),
};

// Exercise API services
export const ExerciseService = {
  getAll: () => fetchApi<Exercise[]>('/exercises'),
  
  getById: (id: number) => fetchApi<Exercise>(`/exercises/${id}`),
  
  create: (exercise: Exercise) => fetchApi<Exercise>('/exercises', {
    method: 'POST',
    body: JSON.stringify(exercise),
  }),
  
  update: (id: number, exercise: Exercise) => fetchApi<Exercise>(`/exercises/${id}`, {
    method: 'PUT',
    body: JSON.stringify(exercise),
  }),
  
  delete: (id: number) => fetchApi<void>(`/exercises/${id}`, {
    method: 'DELETE',
  }),
};

// Set API services
export const SetService = {
  getAll: (exerciseId: number) => fetchApi<Set[]>(`/exercises/${exerciseId}/sets`),
  
  getById: (id: number) => fetchApi<Set>(`/sets/${id}`),
  
  create: (set: Set) => fetchApi<Set>('/sets', {
    method: 'POST',
    body: JSON.stringify(set),
  }),
  
  update: (id: number, set: Set) => fetchApi<Set>(`/sets/${id}`, {
    method: 'PUT',
    body: JSON.stringify(set),
  }),
  
  delete: (id: number) => fetchApi<void>(`/sets/${id}`, {
    method: 'DELETE',
  }),
};

// SuperSet API services
export const SuperSetService = {
  getAll: () => fetchApi<SuperSet[]>('/supersets'),
  
  getById: (id: number) => fetchApi<SuperSet>(`/supersets/${id}`),
  
  create: (superSet: SuperSet) => fetchApi<SuperSet>('/supersets', {
    method: 'POST',
    body: JSON.stringify(superSet),
  }),
  
  update: (id: number, superSet: SuperSet) => fetchApi<SuperSet>(`/supersets/${id}`, {
    method: 'PUT',
    body: JSON.stringify(superSet),
  }),
  
  delete: (id: number) => fetchApi<void>(`/supersets/${id}`, {
    method: 'DELETE',
  }),
};

// Routine API services
export const RoutineService = {
  getAll: () => fetchApi<Routine[]>('/routines'),
  
  getById: (id: number) => fetchApi<Routine>(`/routines/${id}`),
  
  create: (routine: Routine) => fetchApi<Routine>('/routines', {
    method: 'POST',
    body: JSON.stringify(routine),
  }),
  
  update: (id: number, routine: Routine) => fetchApi<Routine>(`/routines/${id}`, {
    method: 'PUT',
    body: JSON.stringify(routine),
  }),
  
  delete: (id: number) => fetchApi<void>(`/routines/${id}`, {
    method: 'DELETE',
  }),
};

// RecordRoutine (Workout) API services
export const WorkoutService = {
  getAll: () => fetchApi<RecordRoutine[]>('/recordroutines'),
  
  getById: (id: number) => fetchApi<RecordRoutine>(`/recordroutines/${id}`),
  
  create: (workout: RecordRoutine) => fetchApi<RecordRoutine>('/recordroutines', {
    method: 'POST',
    body: JSON.stringify(workout),
  }),
  
  update: (id: number, workout: RecordRoutine) => fetchApi<RecordRoutine>(`/recordroutines/${id}`, {
    method: 'PUT',
    body: JSON.stringify(workout),
  }),
  
  delete: (id: number) => fetchApi<void>(`/recordroutines/${id}`, {
    method: 'DELETE',
  }),
  
  // Additional method to get workout statistics for the home page
  getStats: () => fetchApi<WorkoutStats>('/stats'),
};

// RecordExercise API services
export const RecordExerciseService = {
  getAll: (recordRoutineId: number) => fetchApi<RecordExercise[]>(`/recordroutines/${recordRoutineId}/exercises`),
  
  getById: (id: number) => fetchApi<RecordExercise>(`/recordexercises/${id}`),
  
  create: (recordExercise: RecordExercise) => fetchApi<RecordExercise>('/recordexercises', {
    method: 'POST',
    body: JSON.stringify(recordExercise),
  }),
  
  update: (id: number, recordExercise: RecordExercise) => fetchApi<RecordExercise>(`/recordexercises/${id}`, {
    method: 'PUT',
    body: JSON.stringify(recordExercise),
  }),
  
  delete: (id: number) => fetchApi<void>(`/recordexercises/${id}`, {
    method: 'DELETE',
  }),
};

// User profile service
export const ProfileService = {
  get: async () => {
    const user = await fetchApi<User>('/users/1');
    user.birthDate = new Date(user.birthDate).toISOString();
    return user;
  },
  
  update: async (profile: User) => {
    profile.birthDate = new Date(profile.birthDate).toISOString();
    profile.isFemale = profile.isFemale === 'true';
    profile.weight = +profile.weight;
    profile.height = +profile.height;

    return await fetchApi<User>('/users/1', {
    method: 'PUT',
    body: JSON.stringify(profile),
  })},
};
