import { useEffect, useState } from 'react';
import { FaCalendarCheck, FaClock, FaDumbbell } from 'react-icons/fa';
import type { RecordRoutine, WorkoutStats } from '../types/models';
import { WorkoutService } from '../services/api';

const HomePage = () => {
  const [stats, setStats] = useState<WorkoutStats | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);
        const data = await WorkoutService.getStats();
        setStats(data);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch workout stats:', err);
        setError('Could not load workout statistics. Please try again later.');
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return new Intl.DateTimeFormat('en-US', {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  };

  const calculateDuration = (start: string, end: string) => {
    const startTime = new Date(start).getTime();
    const endTime = new Date(end).getTime();
    const durationMinutes = Math.round((endTime - startTime) / (1000 * 60));
    
    return `${durationMinutes} min`;
  };

  if (loading) {
    return (
      <div className="page home-page">
        <h1>Home</h1>
        <div className="loading">Loading workout data...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="page home-page">
        <h1>Home</h1>
        <div className="error-message">{error}</div>
      </div>
    );
  }

  // Display placeholder if no stats
  if (!stats) {
    return (
      <div className="page home-page">
        <h1>Workout Overview</h1>
        <div className="card">
          <h2>Welcome to Go Lift!</h2>
          <p>Start by adding exercises and creating your first workout routine.</p>
          <div className="mt-lg">
            <a href="/workouts" className="btn btn-primary">Go to Workouts</a>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="page home-page">
      <h1>Workout Overview</h1>
      
      {/* Statistics Cards */}
      <div className="stats-grid">
        <div className="card stat-card">
          <div className="stat-icon">
            <FaCalendarCheck size={24} />
          </div>
          <div className="stat-content">
            <div className="stat-value">{stats.totalWorkouts}</div>
            <div className="stat-label">Total Workouts</div>
          </div>
        </div>
        
        <div className="card stat-card">
          <div className="stat-icon">
            <FaClock size={24} />
          </div>
          <div className="stat-content">
            <div className="stat-value">{stats.totalMinutes}</div>
            <div className="stat-label">Total Minutes</div>
          </div>
        </div>
        
        <div className="card stat-card">
          <div className="stat-icon">
            <FaDumbbell size={24} />
          </div>
          <div className="stat-content">
            <div className="stat-value">{stats.totalExercises}</div>
            <div className="stat-label">Exercises Done</div>
          </div>
        </div>
      </div>

      {/* Favorite Data */}
      {(stats.mostFrequentExercise || stats.mostFrequentRoutine) && (
        <div className="card mb-lg">
          <h2>Your Favorites</h2>
          {stats.mostFrequentRoutine && (
            <div className="favorite-item">
              <div className="favorite-label">Most Used Routine:</div>
              <div className="favorite-value">{stats.mostFrequentRoutine.name} ({stats.mostFrequentRoutine.count}x)</div>
            </div>
          )}
          {stats.mostFrequentExercise && (
            <div className="favorite-item">
              <div className="favorite-label">Most Performed Exercise:</div>
              <div className="favorite-value">{stats.mostFrequentExercise.name} ({stats.mostFrequentExercise.count}x)</div>
            </div>
          )}
        </div>
      )}

      {/* Recent Workouts */}
      <h2>Recent Workouts</h2>
      {stats.recentWorkouts && stats.recentWorkouts.length > 0 ? (
        stats.recentWorkouts.map((workout: RecordRoutine) => (
          <div key={workout.id} className="card workout-card">
            <div className="workout-header">
              <h3>{workout.routine?.name || 'Workout'}</h3>
              <div className="workout-date">{formatDate(workout.startedAt)}</div>
            </div>
            {workout.endedAt && (
              <div className="workout-duration">
                Duration: {calculateDuration(workout.startedAt, workout.endedAt)}
              </div>
            )}
            <div className="workout-exercises">
            </div>
          </div>
        ))
      ) : (
        <div className="empty-state">
          <p>No workouts recorded yet. Start your fitness journey today!</p>
          <a href="/new-workout" className="btn btn-primary mt-md">Record a Workout</a>
        </div>
      )}
      
      <style>{`
        .stats-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
          gap: var(--spacing-md);
          margin-bottom: var(--spacing-lg);
        }
        
        .stat-card {
          display: flex;
          align-items: center;
        }
        
        .stat-icon {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 50px;
          height: 50px;
          background-color: rgba(0, 122, 255, 0.1);
          border-radius: 50%;
          color: var(--primary-color);
          margin-right: var(--spacing-md);
        }
        
        .stat-value {
          font-size: 24px;
          font-weight: bold;
        }
        
        .stat-label {
          color: var(--dark-gray);
        }
        
        .favorite-item {
          display: flex;
          justify-content: space-between;
          padding: var(--spacing-sm) 0;
          border-bottom: 1px solid var(--light-gray);
        }
        
        .favorite-item:last-child {
          border-bottom: none;
        }
        
        .favorite-label {
          color: var(--dark-gray);
        }
        
        .workout-card {
          margin-bottom: var(--spacing-md);
        }
        
        .workout-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: var(--spacing-sm);
        }
        
        .workout-date {
          color: var(--dark-gray);
          font-size: 0.9rem;
        }
        
        .workout-duration, .workout-exercises, .workout-notes {
          margin-bottom: var(--spacing-sm);
        }
        
        .workout-feeling {
          color: var(--warning-color);
        }
        
        .empty-state {
          text-align: center;
          padding: var(--spacing-xl);
        }
      `}</style>
    </div>
  );
};

export default HomePage;
