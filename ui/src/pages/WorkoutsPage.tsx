import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { FaPlus, FaPlay, FaEdit, FaTrash, FaFilter } from 'react-icons/fa';
import type { Routine } from '../types/models';
import { RoutineService } from '../services/api';

const WorkoutsPage = () => {
  const [routines, setRoutines] = useState<Routine[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState<string>('');
  const navigate = useNavigate();

  useEffect(() => {
    const fetchRoutines = async () => {
      try {
        setLoading(true);
        const data = await RoutineService.getAll();
        setRoutines(data);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch routines:', err);
        setError('Could not load workout routines. Please try again later.');
      } finally {
        setLoading(false);
      }
    };

    fetchRoutines();
  }, []);

  const handleStartWorkout = (routine: Routine) => {
    // Navigate to new workout page with the selected routine
    navigate('/new-workout', { state: { routineId: routine.id } });
  };

  const handleDeleteRoutine = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this routine?')) {
      try {
        await RoutineService.delete(id);
        setRoutines(routines.filter(routine => routine.id !== id));
      } catch (err) {
        console.error('Failed to delete routine:', err);
        alert('Failed to delete routine. Please try again.');
      }
    }
  };

  const filteredRoutines = routines.filter(routine => 
    routine.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    routine.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

  if (loading) {
    return (
      <div className="page workouts-page">
        <h1>Workouts</h1>
        <div className="loading">Loading routines...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="page workouts-page">
        <h1>Workouts</h1>
        <div className="error-message">{error}</div>
        <div className="mt-lg">
          <Link to="/new-routine" className="btn btn-primary">
            <FaPlus /> Create New Routine
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="page workouts-page">
      <h1>Workout Routines</h1>
      
      {/* Search and Filter */}
      <div className="search-bar">
        <div className="search-input-container">
          <FaFilter />
          <input 
            type="text" 
            placeholder="Search routines..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="search-input"
          />
        </div>
      </div>
      
      {/* Create New Button */}
      <div className="action-buttons mb-lg">
        <Link to="/new-routine" className="btn btn-primary">
          <FaPlus /> Create New Routine
        </Link>
      </div>

      {/* Routines List */}
      {filteredRoutines.length > 0 ? (
        <div className="routines-list">
          {filteredRoutines.map(routine => (
            <div key={routine.id} className="card routine-card">
              <div className="routine-info">
                <h3>{routine.name}</h3>
                <p className="routine-description">{routine.description}</p>
                <div className="routine-stats"></div>
              </div>
              
              <div className="routine-actions">
                <button 
                  className="btn btn-primary"
                  onClick={() => handleStartWorkout(routine)}
                >
                  <FaPlay /> Start
                </button>
                <div className="routine-action-buttons">
                  <Link 
                    to={`/new-routine`} 
                    state={{ editRoutine: routine }}
                    className="btn btn-secondary action-btn"
                  >
                    <FaEdit />
                  </Link>
                  <button 
                    className="btn btn-danger action-btn"
                    onClick={() => routine.id && handleDeleteRoutine(routine.id)}
                  >
                    <FaTrash />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="empty-state">
          {searchTerm ? (
            <p>No routines found matching "{searchTerm}"</p>
          ) : (
            <>
              <p>You haven't created any workout routines yet.</p>
              <p className="mt-sm">Create your first routine to get started!</p>
              <Link to="/new-routine" className="btn btn-primary mt-md">
                <FaPlus /> Create Routine
              </Link>
            </>
          )}
        </div>
      )}

      <style>{`
        .search-bar {
          margin-bottom: var(--spacing-md);
        }
        
        .search-input-container {
          display: flex;
          align-items: center;
          background-color: white;
          border-radius: var(--border-radius);
          padding: 0 var(--spacing-md);
          border: 1px solid var(--light-gray);
        }
        
        .search-input {
          border: none;
          padding: var(--spacing-sm) var(--spacing-sm);
          flex: 1;
        }
        
        .search-input:focus {
          outline: none;
        }
        
        .action-buttons {
          display: flex;
          justify-content: flex-end;
          margin: var(--spacing-md) 0;
        }
        
        .routines-list {
          display: grid;
          gap: var(--spacing-md);
        }
        
        .routine-card {
          display: flex;
          flex-direction: column;
        }
        
        .routine-info {
          flex: 1;
          margin-bottom: var(--spacing-md);
        }
        
        .routine-description {
          color: var(--dark-gray);
          margin: var(--spacing-sm) 0;
        }
        
        .routine-stats {
          display: flex;
          gap: var(--spacing-md);
          color: var(--dark-gray);
          font-size: 0.9rem;
        }
        
        .routine-actions {
          display: flex;
          justify-content: space-between;
          align-items: center;
        }
        
        .routine-action-buttons {
          display: flex;
          gap: var(--spacing-sm);
        }
        
        .action-btn {
          padding: 8px;
          min-width: 40px;
        }
        
        .empty-state {
          text-align: center;
          padding: var(--spacing-xl) var(--spacing-md);
          background-color: white;
          border-radius: var(--border-radius);
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
        
        @media (min-width: 768px) {
          .routine-card {
            flex-direction: row;
          }
          
          .routine-info {
            margin-bottom: 0;
            margin-right: var(--spacing-lg);
          }
          
          .routine-actions {
            flex-direction: column;
            align-items: flex-end;
            justify-content: space-between;
          }
        }
      `}</style>
    </div>
  );
};

export default WorkoutsPage;
