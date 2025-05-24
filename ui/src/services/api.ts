import type { User, Exercise, Muscle, Routine, RecordRoutine, WorkoutStats } from '../types/models';

const API_BASE = '/api';

class BaseService<T> {
  protected endpoint: string;
  constructor(endpoint: string) {
    this.endpoint = endpoint;
  }

  protected async request<R>(path: string, options?: RequestInit): Promise<R> {
    const response = await fetch(`${API_BASE}${path}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    return response.json();
  }

  async getAll(): Promise<T[]> {
    return this.request<T[]>(this.endpoint);
  }

  async get(id: number): Promise<T> {
    return this.request<T>(`${this.endpoint}/${id}`);
  }

  async create(data: Omit<T, 'id' | 'createdAt' | 'updatedAt'>): Promise<T> {
    return this.request<T>(this.endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async update(id: number, data: Partial<T>): Promise<T> {
    return this.request<T>(`${this.endpoint}/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete(id: number): Promise<void> {
    await this.request<void>(`${this.endpoint}/${id}`, {
      method: 'DELETE',
    });
  }
}

class UserService extends BaseService<User> {
  constructor() {
    super('/users');
  }

  // Override get to default to user ID 1
  async get(id: number = 1): Promise<User> {
    return super.get(id);
  }
}

class ExerciseService extends BaseService<Exercise> {
  constructor() {
    super('/exercises');
  }
}

class MuscleService extends BaseService<Muscle> {
  constructor() {
    super('/muscles');
  }
}

class RoutineService extends BaseService<Routine> {
  constructor() {
    super('/routines');
  }
}

class RecordService extends BaseService<RecordRoutine> {
  constructor() {
    super('/records');
  }
}

class StatsService {
  protected async request<R>(path: string, options?: RequestInit): Promise<R> {
    const response = await fetch(`${API_BASE}${path}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    return response.json();
  }

  async get(): Promise<WorkoutStats> {
    return this.request<WorkoutStats>('/stats');
  }
}

class HealthService {
  protected async request<R>(path: string, options?: RequestInit): Promise<R> {
    const response = await fetch(`${API_BASE}${path}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    return response.json();
  }

  async ping(): Promise<{ message: string }> {
    return this.request<{ message: string }>('/ping');
  }
}

// Export service instances
export const userService = new UserService();
export const exerciseService = new ExerciseService();
export const muscleService = new MuscleService();
export const routineService = new RoutineService();
export const recordService = new RecordService();
export const statsService = new StatsService();
export const healthService = new HealthService();
