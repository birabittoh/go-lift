import { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { 
  FaArrowLeft, 
  FaCheck, 
  FaSave, 
  FaPlay,
  FaStop,
  FaStar,
  FaRegStar,
  FaForward
} from 'react-icons/fa';
import type { 
  Routine, 
  Exercise, 
  RecordRoutine,
  RecordSet,
  Set,
  RecordItem
} from '../types/models';
import { routineService } from '../services/api';

interface SetForWorkout {
  id?: number;
  setId: number;
  actualReps: number;
  actualWeight: number;
  actualDuration: number;
  originalSet: Set;
  completed: boolean;
}

interface ExerciseForWorkout {
  id?: number;
  exerciseId: number;
  exercise: Exercise;
  sets: SetForWorkout[];
  startedAt?: string;
  endedAt?: string;
  actualRestTime: number; // in seconds
  notes: string;
}

const NewWorkoutPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  
  // Routines state
  const [routines, setRoutines] = useState<Routine[]>([]);
  const [selectedRoutine, setSelectedRoutine] = useState<Routine | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  
  // Workout tracking state
  const [workoutStarted, setWorkoutStarted] = useState<boolean>(false);
  const [workoutCompleted, setWorkoutCompleted] = useState<boolean>(false);
  const [startTime, setStartTime] = useState<string>('');
  const [endTime, setEndTime] = useState<string | null>(null);
  const [elapsedSeconds, setElapsedSeconds] = useState<number>(0);
  const [intervalId, setIntervalId] = useState<number | null>(null);
  
  // Exercise tracking state
  const [currentExerciseIndex, setCurrentExerciseIndex] = useState<number>(0);
  const [workoutExercises, setWorkoutExercises] = useState<ExerciseForWorkout[]>([]);
  
  // Workout notes and rating
  const [workoutNotes, setWorkoutNotes] = useState<string>('');
  const [feelingRating, setFeelingRating] = useState<number>(3);
  
  // Success message state
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  
  // Load routines and check for pre-selected routine
  useEffect(() => {
    const fetchRoutines = async () => {
      try {
        setIsLoading(true);
        const data = await routineService.getAll();
        setRoutines(data);
        
        // Check if a routine was pre-selected (from workouts page)
        if (location.state && location.state.routineId) {
          const routineId = location.state.routineId;
          const routine = data.find(r => r.id === routineId);
          
          if (routine) {
            handleSelectRoutine(routine);
          }
        }
        
        setError(null);
      } catch (err) {
        console.error('Failed to fetch routines:', err);
        setError('Could not load workout routines. Please try again later.');
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchRoutines();
  }, [location]);
  
  // Setup the workout when a routine is selected
  const handleSelectRoutine = (routine: Routine) => {
    setSelectedRoutine(routine);
    
    // Initialize workout exercises from routine items
    const exercises: ExerciseForWorkout[] = [];
    
    // Process routine items into exercises for the workout
    routine.items.forEach(item => {
      if (item.exerciseItems && item.exerciseId) {
        // This is a regular exercise item
        const exercise = item.exercise;
        
        // Get the sets from the exercise or create default ones
        const exerciseSets = exercise.sets || [];
        const setsForWorkout: SetForWorkout[] = exerciseSets.map(set => ({
          setId: set.id || 0,
          originalSet: set,
          actualReps: set.reps,
          actualWeight: set.weight,
          actualDuration: set.duration,
          completed: false
        }));
        
        // If there are no sets defined, create a default set
        if (setsForWorkout.length === 0) {
          setsForWorkout.push({
            setId: 0,
            originalSet: {
              id: 0,
              exerciseId: exercise.id || 0,
              reps: 10,
              weight: 0,
              duration: 0,
              orderIndex: 0
            },
            actualReps: 10,
            actualWeight: 0,
            actualDuration: 0,
            completed: false
          });
        }
        
        exercises.push({
          exerciseId: exercise.id || 0,
          exercise: exercise,
          sets: setsForWorkout,
          actualRestTime: item.restTime,
          notes: ''
        });
      }
      // We could handle supersets here if needed
    });
    
    setWorkoutExercises(exercises);
    setCurrentExerciseIndex(0);
  };
  
  // Start the workout
  const startWorkout = () => {
    if (!selectedRoutine) return;
    
    const now = new Date().toISOString();
    setStartTime(now);
    setWorkoutStarted(true);
    
    // Mark first exercise as started
    if (workoutExercises.length > 0) {
      const updatedExercises = [...workoutExercises];
      updatedExercises[0].startedAt = now;
      setWorkoutExercises(updatedExercises);
    }
    
    // Start the timer
    const id = window.setInterval(() => {
      setElapsedSeconds(prev => prev + 1);
    }, 1000);
    
    setIntervalId(id);
  };
  
  // Complete the workout
  const completeWorkout = () => {
    if (intervalId) {
      clearInterval(intervalId);
      setIntervalId(null);
    }
    
    const now = new Date().toISOString();
    setEndTime(now);
    
    // Mark current exercise as completed if not already
    if (workoutExercises.length > 0 && currentExerciseIndex < workoutExercises.length) {
      const updatedExercises = [...workoutExercises];
      const currentExercise = updatedExercises[currentExerciseIndex];
      
      if (currentExercise.startedAt && !currentExercise.endedAt) {
        currentExercise.endedAt = now;
      }
      
      setWorkoutExercises(updatedExercises);
    }
    
    setWorkoutCompleted(true);
  };
  
  // Format timer display
  const formatTime = (seconds: number) => {
    const hrs = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    
    return `${hrs > 0 ? hrs + ':' : ''}${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };
  
  // Handle set completion toggle
  const toggleSetCompleted = (exerciseIndex: number, setIndex: number) => {
    const updatedExercises = [...workoutExercises];
    const currentSet = updatedExercises[exerciseIndex].sets[setIndex];
    currentSet.completed = !currentSet.completed;
    
    setWorkoutExercises(updatedExercises);
  };
  
  // Handle weight, reps or duration change
  const handleSetDataChange = (
    exerciseIndex: number,
    setIndex: number,
    field: 'actualReps' | 'actualWeight' | 'actualDuration',
    value: number
  ) => {
    const updatedExercises = [...workoutExercises];
    const currentSet = updatedExercises[exerciseIndex].sets[setIndex];
    currentSet[field] = value;
    
    setWorkoutExercises(updatedExercises);
  };
  
  // Move to next exercise
  const nextExercise = () => {
    if (currentExerciseIndex >= workoutExercises.length - 1) return;
    
    const now = new Date().toISOString();
    const updatedExercises = [...workoutExercises];
    
    // Complete current exercise
    const currentExercise = updatedExercises[currentExerciseIndex];
    if (currentExercise.startedAt && !currentExercise.endedAt) {
      currentExercise.endedAt = now;
    }
    
    // Start next exercise
    const nextIndex = currentExerciseIndex + 1;
    const nextExercise = updatedExercises[nextIndex];
    nextExercise.startedAt = now;
    
    setWorkoutExercises(updatedExercises);
    setCurrentExerciseIndex(nextIndex);
  };
  
  // Handle notes for an exercise
  const handleExerciseNotes = (exerciseIndex: number, notes: string) => {
    const updatedExercises = [...workoutExercises];
    updatedExercises[exerciseIndex].notes = notes;
    
    setWorkoutExercises(updatedExercises);
  };
  
  // Create RecordSets from workout exercise sets
  const createRecordSets = (exercise: ExerciseForWorkout): RecordSet[] => {
    return exercise.sets.map((set, index) => ({
      recordExerciseId: 0, // Will be filled in by backend
      setId: set.setId,
      actualReps: set.actualReps,
      actualWeight: set.actualWeight,
      actualDuration: set.actualDuration,
      completedAt: exercise.endedAt || new Date().toISOString(),
      orderIndex: index,
      set: set.originalSet
    }));
  };
  
  // Save workout record
  const saveWorkout = async () => {
    if (!selectedRoutine || !startTime) return;
    
    try {
      const now = new Date().toISOString();
      
      // Ensure all exercises have start/end times
      const completedExercises = workoutExercises.map((ex) => {
        if (!ex.startedAt) {
          ex.startedAt = startTime;
        }
        if (!ex.endedAt) {
          ex.endedAt = endTime || now;
        }
        return ex;
      });
      
      // Create RecordExercises from completed exercises
      const recordExercises: RecordExercise[] = completedExercises.map((ex, index) => ({
        id: undefined,
        recordRoutineId: 0, // Will be filled in by backend
        exerciseId: ex.exerciseId,
        startedAt: ex.startedAt || startTime,
        endedAt: ex.endedAt || now,
        actualRestTime: ex.actualRestTime,
        orderIndex: index,
        recordSets: createRecordSets(ex),
        exercise: ex.exercise
      }));
      
      // Create RecordRoutineItems from recordExercises
      const recordItems: RecordItem[] = recordExercises.map((ex, index) => ({
        recordRoutineId: 0, // Will be filled in by backend
        recordExerciseId: undefined, // Will be filled in after recordExercise is created
        recordSuperSetId: null,
        actualRestTime: workoutExercises[index].actualRestTime,
        orderIndex: index,
        recordExercise: ex,
        recordSuperSet: null
      }));
      
      const workoutRecord: RecordRoutine = {
        routineId: selectedRoutine.id!,
        duration: elapsedSeconds,
        routine: selectedRoutine,
        recordItems: recordItems
      };
      
      await WorkoutService.create(workoutRecord);
      setSuccessMessage('Workout saved successfully!');
      
      // Redirect after a brief delay
      setTimeout(() => {
        navigate('/home');
      }, 1500);
    } catch (err) {
      console.error('Failed to save workout:', err);
      setError('Failed to save your workout. Please try again.');
    }
  };
  
  // Check if all sets in current exercise are completed
  const isCurrentExerciseComplete = () => {
    if (currentExerciseIndex >= workoutExercises.length) return false;
    
    const currentExercise = workoutExercises[currentExerciseIndex];
    return currentExercise.sets.every(set => set.completed);
  };
  
  // Progress status percentage
  const calculateProgress = () => {
    if (workoutExercises.length === 0) return 0;
    
    const totalSets = workoutExercises.reduce((total, ex) => total + ex.sets.length, 0);
    const completedSets = workoutExercises.reduce((total, ex) => {
      return total + ex.sets.filter(set => set.completed).length;
    }, 0);
    
    return Math.round((completedSets / totalSets) * 100);
  };
  
  return (
    <div className="page new-workout-page">
      <div className="page-header">
        <button 
          onClick={() => navigate(-1)} 
          className="btn btn-secondary back-button"
        >
          <FaArrowLeft /> Back
        </button>
        <h1>New Workout</h1>
      </div>
      
      {error && <div className="error-message">{error}</div>}
      {successMessage && <div className="success-message">{successMessage}</div>}
      
      {!selectedRoutine ? (
        // Routine selection view
        <div className="routine-selection card">
          <h2>Select a Routine</h2>
          
          {isLoading ? (
            <div className="loading">Loading routines...</div>
          ) : routines.length === 0 ? (
            <div className="empty-state">
              <p>No routines found.</p>
              <p>Create a routine to start working out.</p>
              <button 
                onClick={() => navigate('/new-routine')}
                className="btn btn-primary mt-md"
              >
                Create Routine
              </button>
            </div>
          ) : (
            <div className="routines-list">
              {routines.map(routine => (
                <div key={routine.id} className="routine-item" onClick={() => handleSelectRoutine(routine)}>
                  <h3>{routine.name}</h3>
                  {routine.description && <p>{routine.description}</p>}
                  <div className="routine-meta">
                    <span>{routine.routineItems.length} exercises</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      ) : !workoutStarted ? (
        // Workout ready view
        <div className="workout-ready card">
          <h2>{selectedRoutine.name}</h2>
          {selectedRoutine.description && <p className="routine-description">{selectedRoutine.description}</p>}
          
          <div className="workout-details">
            <div className="detail-item">
              <span className="detail-label">Exercises:</span>
              <span className="detail-value">{workoutExercises.length}</span>
            </div>
            <div className="detail-item">
              <span className="detail-label">Sets:</span>
              <span className="detail-value">
                {workoutExercises.reduce((total, ex) => total + ex.sets.length, 0)}
              </span>
            </div>
          </div>
          
          <div className="exercise-preview">
            <h3>Exercises</h3>
            <ul className="exercise-list">
              {workoutExercises.map((exercise, index) => (
                <li key={`${exercise.exerciseId}-${index}`}>
                  <div className="exercise-name">{exercise.exercise.name}</div>
                  <div className="exercise-sets">{exercise.sets.length} sets</div>
                </li>
              ))}
            </ul>
          </div>
          
          <div className="action-buttons">
            <button 
              className="btn btn-primary btn-lg btn-block"
              onClick={startWorkout}
            >
              <FaPlay /> Start Workout
            </button>
            <button 
              className="btn btn-secondary btn-block mt-md"
              onClick={() => setSelectedRoutine(null)}
            >
              Select Different Routine
            </button>
          </div>
        </div>
      ) : workoutCompleted ? (
        // Workout complete view
        <div className="workout-complete card">
          <div className="workout-summary">
            <h2>Workout Complete!</h2>
            
            <div className="summary-stats">
              <div className="stat">
                <span className="stat-label">Duration</span>
                <span className="stat-value">{formatTime(elapsedSeconds)}</span>
              </div>
              
              <div className="stat">
                <span className="stat-label">Completed</span>
                <span className="stat-value">{calculateProgress()}%</span>
              </div>
            </div>
            
            <div className="feeling-rating">
              <p>How was your workout?</p>
              <div className="stars">
                {[1, 2, 3, 4, 5].map(rating => (
                  <button
                    key={rating}
                    onClick={() => setFeelingRating(rating)}
                    className="star-btn"
                  >
                    {rating <= feelingRating ? <FaStar /> : <FaRegStar />}
                  </button>
                ))}
              </div>
            </div>
            
            <div className="workout-notes form-group">
              <label htmlFor="workout-notes">Workout Notes</label>
              <textarea
                id="workout-notes"
                value={workoutNotes}
                onChange={e => setWorkoutNotes(e.target.value)}
                placeholder="Add notes about the overall workout..."
                rows={3}
              />
            </div>
            
            <div className="action-buttons">
              <button
                onClick={saveWorkout}
                className="btn btn-primary btn-lg btn-block"
              >
                <FaSave /> Save Workout
              </button>
            </div>
          </div>
        </div>
      ) : (
        // Active workout view
        <div className="active-workout">
          <div className="workout-header card">
            <h2>{selectedRoutine.name}</h2>
            
            <div className="workout-timer">
              <div className="timer-value">{formatTime(elapsedSeconds)}</div>
              <div className="progress-bar">
                <div 
                  className="progress" 
                  style={{ width: `${calculateProgress()}%` }}
                ></div>
              </div>
            </div>
          </div>
          
          <div className="current-exercise card">
            {currentExerciseIndex < workoutExercises.length ? (
              <>
                <h3 className="exercise-name">
                  {workoutExercises[currentExerciseIndex].exercise.name}
                </h3>
                
                <div className="exercise-sets">
                  <table className="sets-table">
                    <thead>
                      <tr>
                        <th>Set</th>
                        <th>Weight</th>
                        <th>Reps</th>
                        {workoutExercises[currentExerciseIndex].sets.some(s => s.originalSet.duration > 0) && (
                          <th>Time</th>
                        )}
                        <th>Done</th>
                      </tr>
                    </thead>
                    <tbody>
                      {workoutExercises[currentExerciseIndex].sets.map((set, setIndex) => (
                        <tr key={setIndex} className={set.completed ? 'completed' : ''}>
                          <td>{setIndex + 1}</td>
                          <td>
                            <input
                              type="number"
                              min="0"
                              step="1"
                              value={set.actualWeight}
                              onChange={e => handleSetDataChange(
                                currentExerciseIndex,
                                setIndex,
                                'actualWeight',
                                parseFloat(e.target.value) || 0
                              )}
                            />
                          </td>
                          <td>
                            <input
                              type="number"
                              min="0"
                              step="1"
                              value={set.actualReps}
                              onChange={e => handleSetDataChange(
                                currentExerciseIndex,
                                setIndex,
                                'actualReps',
                                parseInt(e.target.value) || 0
                              )}
                            />
                          </td>
                          {workoutExercises[currentExerciseIndex].sets.some(s => s.originalSet.duration > 0) && (
                            <td>
                              <input
                                type="number"
                                min="0"
                                step="1"
                                value={set.actualDuration}
                                onChange={e => handleSetDataChange(
                                  currentExerciseIndex,
                                  setIndex,
                                  'actualDuration',
                                  parseInt(e.target.value) || 0
                                )}
                              />
                            </td>
                          )}
                          <td>
                            <button
                              className={`btn-check ${set.completed ? 'completed' : ''}`}
                              onClick={() => toggleSetCompleted(currentExerciseIndex, setIndex)}
                            >
                              {set.completed && <FaCheck />}
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                
                <div className="exercise-notes form-group">
                  <label htmlFor="exercise-notes">Exercise Notes</label>
                  <textarea
                    id="exercise-notes"
                    value={workoutExercises[currentExerciseIndex].notes}
                    onChange={e => handleExerciseNotes(currentExerciseIndex, e.target.value)}
                    placeholder="Add notes for this exercise..."
                    rows={2}
                  />
                </div>
                
                <div className="exercise-navigation">
                  {currentExerciseIndex < workoutExercises.length - 1 && (
                    <button
                      onClick={nextExercise}
                      className={`btn btn-primary ${isCurrentExerciseComplete() ? 'pulse' : ''}`}
                      disabled={!isCurrentExerciseComplete() && workoutExercises[currentExerciseIndex].sets.length > 0}
                    >
                      <FaForward /> Next Exercise
                    </button>
                  )}
                  
                  {currentExerciseIndex === workoutExercises.length - 1 && isCurrentExerciseComplete() && (
                    <button
                      onClick={completeWorkout}
                      className="btn btn-primary pulse"
                    >
                      <FaStop /> Finish Workout
                    </button>
                  )}
                </div>
              </>
            ) : (
              <div>
                <p>All exercises completed!</p>
                <button
                  onClick={completeWorkout}
                  className="btn btn-primary"
                >
                  <FaStop /> Finish Workout
                </button>
              </div>
            )}
          </div>
          
          <div className="workout-nav">
            <div className="exercises-list card">
              <h3>Progress</h3>
              <ul>
                {workoutExercises.map((ex, index) => (
                  <li 
                    key={`${ex.exerciseId}-${index}`} 
                    className={`
                      ${index === currentExerciseIndex ? 'active' : ''}
                      ${ex.sets.every(s => s.completed) ? 'completed' : ''}
                    `}
                    onClick={() => setCurrentExerciseIndex(index)}
                  >
                    <span className="exercise-number">{index + 1}</span>
                    <span className="exercise-list-name">{ex.exercise.name}</span>
                    <span className="exercise-progress">
                      {ex.sets.filter(s => s.completed).length}/{ex.sets.length}
                    </span>
                  </li>
                ))}
              </ul>
            </div>
            
            <div className="workout-actions card">
              <button
                onClick={completeWorkout}
                className="btn btn-danger btn-block"
              >
                <FaStop /> End Workout
              </button>
            </div>
          </div>
        </div>
      )}
      
      <style>{`
        .page-header {
          display: flex;
          align-items: center;
          margin-bottom: var(--spacing-lg);
        }
        
        .back-button {
          margin-right: var(--spacing-md);
          padding: var(--spacing-sm) var(--spacing-md);
        }
        
        /* Routine Selection */
        .routines-list {
          display: grid;
          gap: var(--spacing-md);
        }
        
        .routine-item {
          padding: var(--spacing-md);
          border: 1px solid var(--border-color);
          border-radius: var(--border-radius);
          cursor: pointer;
          transition: all 0.2s;
        }
        
        .routine-item:hover {
          background-color: rgba(0, 122, 255, 0.05);
          border-color: var(--primary-color);
        }
        
        .routine-item h3 {
          margin: 0;
          margin-bottom: var(--spacing-xs);
        }
        
        .routine-item p {
          margin: 0;
          margin-bottom: var(--spacing-sm);
          color: var(--text-muted);
        }
        
        .routine-meta {
          font-size: 0.9rem;
          color: var(--text-muted);
        }
        
        /* Workout Ready */
        .routine-description {
          color: var(--text-muted);
          margin-bottom: var(--spacing-lg);
        }
        
        .workout-details {
          display: flex;
          gap: var(--spacing-lg);
          margin-bottom: var(--spacing-lg);
        }
        
        .detail-item {
          display: flex;
          flex-direction: column;
          align-items: center;
          background-color: rgba(0, 122, 255, 0.1);
          padding: var(--spacing-md) var(--spacing-lg);
          border-radius: var(--border-radius);
        }
        
        .detail-label {
          color: var(--text-muted);
          font-size: 0.9rem;
        }
        
        .detail-value {
          font-size: 1.2rem;
          font-weight: bold;
        }
        
        .exercise-preview {
          margin-bottom: var(--spacing-lg);
        }
        
        .exercise-list {
          list-style: none;
          padding: 0;
          margin: 0;
        }
        
        .exercise-list li {
          padding: var(--spacing-sm) 0;
          border-bottom: 1px solid var(--border-color);
          display: flex;
          justify-content: space-between;
        }
        
        .exercise-list li:last-child {
          border-bottom: none;
        }
        
        .exercise-sets {
          color: var(--text-muted);
          font-size: 0.9rem;
        }
        
        .btn-lg {
          padding: var(--spacing-md);
          font-size: 1.1rem;
        }
        
        /* Active Workout */
        .workout-header {
          margin-bottom: var(--spacing-md);
          padding: var(--spacing-md);
        }
        
        .workout-header h2 {
          margin-bottom: var(--spacing-sm);
        }
        
        .workout-timer {
          text-align: center;
        }
        
        .timer-value {
          font-size: 1.5rem;
          font-weight: bold;
          margin-bottom: var(--spacing-xs);
        }
        
        .progress-bar {
          height: 8px;
          background-color: var(--light-gray);
          border-radius: 4px;
          overflow: hidden;
        }
        
        .progress {
          height: 100%;
          background-color: var(--primary-color);
          transition: width 0.3s ease;
        }
        
        .current-exercise {
          margin-bottom: var(--spacing-md);
          padding: var(--spacing-md);
        }
        
        .exercise-name {
          margin-bottom: var(--spacing-md);
          font-size: 1.2rem;
        }
        
        /* Sets table */
        .sets-table {
          width: 100%;
          border-collapse: collapse;
          margin-bottom: var(--spacing-md);
        }
        
        .sets-table th {
          padding: var(--spacing-sm);
          text-align: center;
          border-bottom: 1px solid var(--border-color);
          font-weight: 600;
        }
        
        .sets-table td {
          padding: var(--spacing-sm);
          text-align: center;
          border-bottom: 1px solid var(--border-color);
        }
        
        .sets-table tr.completed {
          background-color: rgba(52, 199, 89, 0.1);
        }
        
        .sets-table input {
          width: 60px;
          padding: 4px;
          text-align: center;
        }
        
        .btn-check {
          width: 30px;
          height: 30px;
          border-radius: 50%;
          border: 1px solid var(--border-color);
          background: white;
          display: flex;
          align-items: center;
          justify-content: center;
          cursor: pointer;
        }
        
        .btn-check.completed {
          background-color: var(--success-color);
          color: white;
          border-color: var(--success-color);
        }
        
        .exercise-notes {
          margin-bottom: var(--spacing-md);
        }
        
        .exercise-navigation {
          display: flex;
          justify-content: flex-end;
        }
        
        /* Pulse animation for the next button */
        .pulse {
          animation: pulse 1.5s infinite;
        }
        
        @keyframes pulse {
          0% {
            box-shadow: 0 0 0 0 rgba(0, 122, 255, 0.4);
          }
          70% {
            box-shadow: 0 0 0 10px rgba(0, 122, 255, 0);
          }
          100% {
            box-shadow: 0 0 0 0 rgba(0, 122, 255, 0);
          }
        }
        
        /* Progress list */
        .workout-nav {
          margin-bottom: var(--spacing-lg);
        }
        
        .exercises-list ul {
          list-style: none;
          padding: 0;
          margin: 0;
        }
        
        .exercises-list li {
          display: flex;
          align-items: center;
          padding: var(--spacing-sm);
          border-bottom: 1px solid var(--border-color);
          cursor: pointer;
        }
        
        .exercises-list li.active {
          background-color: rgba(0, 122, 255, 0.1);
        }
        
        .exercises-list li.completed {
          color: var(--text-muted);
        }
        
        .exercise-number {
          width: 24px;
          height: 24px;
          border-radius: 50%;
          background-color: var(--light-gray);
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: var(--spacing-sm);
          font-size: 0.8rem;
        }
        
        .exercises-list li.completed .exercise-number {
          background-color: var(--success-color);
          color: white;
        }
        
        .exercise-list-name {
          flex: 1;
        }
        
        .exercise-progress {
          font-size: 0.8rem;
          color: var(--text-muted);
        }
        
        .workout-actions {
          margin-top: var(--spacing-md);
          padding: var(--spacing-md);
        }
        
        /* Workout Complete */
        .workout-complete {
          text-align: center;
          padding: var(--spacing-lg);
        }
        
        .workout-summary h2 {
          margin-bottom: var(--spacing-lg);
        }
        
        .summary-stats {
          display: flex;
          justify-content: center;
          gap: var(--spacing-xl);
          margin-bottom: var(--spacing-lg);
        }
        
        .stat {
          display: flex;
          flex-direction: column;
        }
        
        .stat-label {
          font-size: 0.9rem;
          color: var(--text-muted);
        }
        
        .stat-value {
          font-size: 1.5rem;
          font-weight: bold;
        }
        
        .feeling-rating {
          margin-bottom: var(--spacing-lg);
        }
        
        .stars {
          display: flex;
          justify-content: center;
          gap: var(--spacing-sm);
        }
        
        .star-btn {
          background: none;
          border: none;
          font-size: 1.5rem;
          color: #ffb700;
          cursor: pointer;
        }
        
        /* Shared */
        .loading {
          text-align: center;
          padding: var(--spacing-lg);
          color: var(--text-muted);
        }
        
        .empty-state {
          text-align: center;
          padding: var(--spacing-lg);
          color: var(--text-muted);
        }
        
        .error-message {
          background-color: rgba(255, 59, 48, 0.1);
          color: var(--danger-color);
          padding: var(--spacing-md);
          border-radius: var(--border-radius);
          margin-bottom: var(--spacing-lg);
        }
        
        .success-message {
          background-color: rgba(52, 199, 89, 0.1);
          color: var(--success-color);
          padding: var(--spacing-md);
          border-radius: var(--border-radius);
          margin-bottom: var(--spacing-lg);
        }
        
        .mt-md {
          margin-top: var(--spacing-md);
        }
        
        .card {
          background-color: white;
          border-radius: var(--border-radius);
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          padding: var(--spacing-lg);
        }
        
        @media (min-width: 768px) {
          .active-workout {
            display: grid;
            grid-template-columns: 2fr 1fr;
            gap: var(--spacing-md);
          }
          
          .workout-header {
            grid-column: 1 / -1;
          }
          
          .workout-nav {
            grid-column: 2;
            grid-row: 2;
          }
        }
      `}</style>
    </div>
  );
};

export default NewWorkoutPage;
