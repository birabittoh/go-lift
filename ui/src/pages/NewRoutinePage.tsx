import { useState, useEffect } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import {
  FaPlus,
  FaSave,
  FaArrowLeft,
  FaTrash,
  FaArrowUp,
  FaArrowDown,
  FaFilter,
} from 'react-icons/fa';
import type { Routine, Exercise, RoutineItem } from '../types/models';
import { RoutineService, ExerciseService } from '../services/api';

const NewRoutinePage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  
  // State for exercises 
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [selectedItems, setSelectedItems] = useState<RoutineItem[]>([]);
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [muscleFilter, setMuscleFilter] = useState<string>('');
  
  // State for routine data
  const [routineName, setRoutineName] = useState<string>('');
  const [routineDescription, setRoutineDescription] = useState<string>('');
  const [estimatedDuration, setEstimatedDuration] = useState<number>(0);
  
  // UI state
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isSaving, setIsSaving] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  
  // Track if we're editing an existing routine
  const [isEditMode, setIsEditMode] = useState<boolean>(false);
  const [routineToEdit, setRoutineToEdit] = useState<Routine | null>(null);
  
  // Fetch available exercises and check if we're in edit mode
  useEffect(() => {
    const fetchExercises = async () => {
      try {
        setIsLoading(true);
        const data = await ExerciseService.getAll();
        setExercises(data);
        
        // Check if we're editing an existing routine
        if (location.state && location.state.editRoutine) {
          const routine = location.state.editRoutine as Routine;
          setRoutineName(routine.name);
          setRoutineDescription(routine.description);
          
          // Set selected items from the routine
          if (routine.routineItems && routine.routineItems.length > 0) {
            setSelectedItems(routine.routineItems);
          }
          
          setRoutineToEdit(routine);
          setIsEditMode(true);
        }
        
        setError(null);
      } catch (err) {
        console.error('Failed to fetch exercises:', err);
        setError('Could not load exercises. Please try again later.');
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchExercises();
  }, [location]);
  
  // Find unique muscle groups for filtering
  const muscleGroups = [...new Set(exercises.flatMap(ex => 
    ex.muscleGroups.map(mg => mg.name)
  ))].sort();
  
  // Filter exercises based on search and muscle filter
  const filteredExercises = exercises.filter(ex => {
    const matchesSearch = ex.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
                          ex.description.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesMuscle = !muscleFilter || 
                         ex.muscleGroups.some(mg => mg.name === muscleFilter);
    return matchesSearch && matchesMuscle;
  });
  
  // Handle adding an exercise to the routine
  const handleAddExercise = (exercise: Exercise) => {
    const newItem: RoutineItem = {
      routineId: routineToEdit?.id || 0,
      exerciseId: exercise.id,
      superSetId: null,
      restTime: 60, // Default rest time 60 seconds
      orderIndex: selectedItems.length,
      exercise: exercise,
      superSet: null,
    };
    
    setSelectedItems([...selectedItems, newItem]);
  };
  
  // Handle removing an exercise from the routine
  const handleRemoveItem = (index: number) => {
    const newSelectedItems = [...selectedItems];
    newSelectedItems.splice(index, 1);
    
    // Update order values
    const reorderedItems = newSelectedItems.map((item, i) => ({
      ...item,
      orderIndex: i,
    }));
    
    setSelectedItems(reorderedItems);
  };
  
  // Handle updating exercise details
  const handleItemChange = (index: number, field: string, value: any) => {
    const newSelectedItems = [...selectedItems];
    // @ts-ignore
    newSelectedItems[index][field] = value;
    setSelectedItems(newSelectedItems);
  };
  
  // Handle moving exercises up/down in the order
  const handleMoveItem = (index: number, direction: 'up' | 'down') => {
    if (
      (direction === 'up' && index === 0) || 
      (direction === 'down' && index === selectedItems.length - 1)
    ) {
      return;
    }
    
    const newSelectedItems = [...selectedItems];
    const swapIndex = direction === 'up' ? index - 1 : index + 1;
    
    // Swap items
    [newSelectedItems[index], newSelectedItems[swapIndex]] = 
    [newSelectedItems[swapIndex], newSelectedItems[index]];
    
    // Update order values
    const reorderedItems = newSelectedItems.map((item, i) => ({
      ...item,
      orderIndex: i,
    }));
    
    setSelectedItems(reorderedItems);
  };
  
  // Calculate estimated duration based on exercises and rest times
  useEffect(() => {
    let totalMinutes = 0;
    
    selectedItems.forEach(item => {
      if (item.exercise) {
        // Estimate time for each set (1 min per set) + rest time
        const setCount = item.exercise.sets?.length || 3; // Default to 3 if no sets defined
        const restTime = item.restTime / 60; // Convert seconds to minutes
        totalMinutes += setCount + restTime;
      } else if (item.superSet) {
        // For supersets, account for both exercises
        const primarySets = item.superSet.primaryExercise?.sets?.length || 3;
        const secondarySets = item.superSet.secondaryExercise?.sets?.length || 3;
        const restTime = item.superSet.restTime / 60;
        totalMinutes += primarySets + secondarySets + restTime;
      }
    });
    
    // Add some buffer time for transitions between exercises
    totalMinutes += selectedItems.length > 0 ? Math.ceil(selectedItems.length / 3) : 0;
    
    setEstimatedDuration(Math.ceil(totalMinutes));
  }, [selectedItems]);
  
  // Handle saving the routine
  const handleSaveRoutine = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (selectedItems.length === 0) {
      setError('Please add at least one exercise to your routine.');
      return;
    }
    
    setIsSaving(true);
    setError(null);
    
    // Prepare the routineItems by removing circular references
    const sanitizedItems = selectedItems.map(item => {
      const { exercise, superSet, ...rest } = item;
      return rest;
    });
    
    const routineData: Routine = {
      name: routineName,
      description: routineDescription,
      routineItems: sanitizedItems,
    };
    
    try {
      if (isEditMode && routineToEdit?.id) {
        // Update existing routine
        await RoutineService.update(routineToEdit.id, routineData);
        setSuccessMessage('Routine updated successfully!');
      } else {
        // Create new routine
        await RoutineService.create(routineData);
        setSuccessMessage('Routine created successfully!');
      }
      
      // Show success message briefly then redirect
      setTimeout(() => {
        navigate('/workouts');
      }, 1500);
    } catch (err) {
      console.error('Failed to save routine:', err);
      setError('Failed to save routine. Please try again.');
    } finally {
      setIsSaving(false);
    }
  };
  
  return (
    <div className="page new-routine-page">
      <div className="page-header">
        <button 
          onClick={() => navigate(-1)} 
          className="btn btn-secondary back-button"
        >
          <FaArrowLeft /> Back
        </button>
        <h1>{isEditMode ? 'Edit Routine' : 'Create Routine'}</h1>
      </div>
      
      {error && <div className="error-message">{error}</div>}
      {successMessage && <div className="success-message">{successMessage}</div>}
      
      {isLoading ? (
        <div className="loading">Loading exercises...</div>
      ) : (
        <div className="routine-builder">
          <div className="routine-details card">
            <h2>Routine Details</h2>
            <form onSubmit={handleSaveRoutine}>
              <div className="form-group">
                <label htmlFor="routineName">Routine Name*</label>
                <input
                  type="text"
                  id="routineName"
                  value={routineName}
                  onChange={(e) => setRoutineName(e.target.value)}
                  placeholder="e.g. Upper Body Strength"
                  required
                  disabled={isSaving}
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="routineDescription">Description</label>
                <textarea
                  id="routineDescription"
                  value={routineDescription}
                  onChange={(e) => setRoutineDescription(e.target.value)}
                  placeholder="Describe your routine..."
                  rows={3}
                  disabled={isSaving}
                />
              </div>
              
              <div className="routine-summary">
                <div className="summary-item">
                  <span className="summary-label">Exercises:</span>
                  <span className="summary-value">{selectedItems.length}</span>
                </div>
                <div className="summary-item">
                  <span className="summary-label">Est. Duration:</span>
                  <span className="summary-value">{estimatedDuration} min</span>
                </div>
              </div>
              
              <div className="selected-exercises">
                <h3>Selected Exercises</h3>
                
                {selectedItems.length === 0 ? (
                  <div className="empty-state">
                    <p>No exercises added yet.</p>
                    <p className="hint">Select exercises from the list below to add them to your routine.</p>
                  </div>
                ) : (
                  <div className="exercise-list">
                    {selectedItems.map((item, index) => (
                      <div key={`${index}-${item.exerciseId || item.superSetId}`} className="selected-exercise-item">
                        <div className="exercise-order">#{index + 1}</div>
                        
                        <div className="exercise-content">
                          <div className="exercise-name">{item.exercise?.name || (item.superSet ? `${item.superSet.primaryExercise.name} + ${item.superSet.secondaryExercise.name}` : "Unknown Exercise")}</div>
                          
                          <div className="exercise-details">
                            <div className="detail-item">
                              <label htmlFor={`rest-${index}`}>Rest (sec):</label>
                              <input
                                id={`rest-${index}`}
                                type="number"
                                min="0"
                                max="300"
                                step="15"
                                value={item.restTime}
                                onChange={(e) => handleItemChange(index, 'restTime', parseInt(e.target.value))}
                                disabled={isSaving}
                              />
                            </div>
                          </div>
                        </div>
                        
                        <div className="exercise-actions">
                          <button
                            type="button"
                            onClick={() => handleMoveItem(index, 'up')}
                            disabled={index === 0 || isSaving}
                            className="btn btn-secondary action-btn"
                            title="Move up"
                          >
                            <FaArrowUp />
                          </button>
                          <button
                            type="button"
                            onClick={() => handleMoveItem(index, 'down')}
                            disabled={index === selectedItems.length - 1 || isSaving}
                            className="btn btn-secondary action-btn"
                            title="Move down"
                          >
                            <FaArrowDown />
                          </button>
                          <button
                            type="button"
                            onClick={() => handleRemoveItem(index)}
                            disabled={isSaving}
                            className="btn btn-danger action-btn"
                            title="Remove"
                          >
                            <FaTrash />
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
              
              <div className="form-actions">
                <button
                  type="submit"
                  className="btn btn-primary btn-block"
                  disabled={isSaving || selectedItems.length === 0 || !routineName}
                >
                  <FaSave /> {isEditMode ? 'Update Routine' : 'Save Routine'}
                </button>
              </div>
            </form>
          </div>
          
          <div className="exercise-picker card">
            <h2>Available Exercises</h2>
            
            <div className="exercise-filters">
              <div className="search-input-container">
                <FaFilter />
                <input
                  type="text"
                  placeholder="Search exercises..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="search-input"
                />
              </div>
              
              <div className="muscle-filter">
                <select
                  value={muscleFilter}
                  onChange={(e) => setMuscleFilter(e.target.value)}
                >
                  <option value="">All Muscle Groups</option>
                  {muscleGroups.map(group => (
                    <option key={group} value={group}>{group}</option>
                  ))}
                </select>
              </div>
            </div>
            
            {filteredExercises.length === 0 && (
              <div className="empty-state">
                <p>No exercises found.</p>
                <Link to="/new-exercise" className="btn btn-primary mt-md">
                  <FaPlus /> Create New Exercise
                </Link>
              </div>
            )}
            
            <div className="available-exercises">
              {filteredExercises.map(exercise => (
                <div key={exercise.id} className="exercise-item">
                  <div className="exercise-info">
                    <h4>{exercise.name}</h4>
                    <div className="exercise-metadata">
                      <span>{exercise.muscleGroups.map(mg => mg.name).join(', ')}</span>
                      <span>{exercise.equipment.map(eq => eq.name).join(', ')}</span>
                    </div>
                    {exercise.description && (
                      <div className="exercise-description">{exercise.description}</div>
                    )}
                  </div>
                  <button
                    type="button"
                    onClick={() => handleAddExercise(exercise)}
                    className="btn btn-primary"
                    disabled={isSaving}
                  >
                    <FaPlus /> Add
                  </button>
                </div>
              ))}
            </div>
            
            <div className="text-center mt-lg">
              <Link to="/new-exercise" className="btn btn-secondary">
                <FaPlus /> Create New Exercise
              </Link>
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
        
        .routine-builder {
          display: grid;
          gap: var(--spacing-lg);
        }
        
        @media (min-width: 768px) {
          .routine-builder {
            grid-template-columns: 1fr 1fr;
          }
        }
        
        .routine-summary {
          display: flex;
          gap: var(--spacing-lg);
          margin: var(--spacing-md) 0;
        }
        
        .summary-item {
          padding: var(--spacing-sm) var(--spacing-md);
          background-color: rgba(0, 122, 255, 0.1);
          border-radius: var(--border-radius);
          display: flex;
          align-items: center;
          gap: var(--spacing-sm);
        }
        
        .summary-label {
          font-weight: 500;
        }
        
        .selected-exercises {
          margin-top: var(--spacing-lg);
        }
        
        h3 {
          margin-bottom: var(--spacing-md);
        }
        
        .selected-exercise-item {
          display: flex;
          align-items: center;
          padding: var(--spacing-md);
          border: 1px solid var(--border-color);
          border-radius: var(--border-radius);
          margin-bottom: var(--spacing-sm);
          background-color: white;
        }
        
        .exercise-order {
          font-weight: bold;
          width: 30px;
          height: 30px;
          border-radius: 50%;
          background-color: var(--primary-color);
          color: white;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: var(--spacing-md);
        }
        
        .exercise-content {
          flex: 1;
        }
        
        .exercise-name {
          font-weight: 600;
          margin-bottom: var(--spacing-xs);
        }
        
        .exercise-details {
          display: flex;
          flex-wrap: wrap;
          gap: var(--spacing-md);
        }
        
        .detail-item {
          display: flex;
          align-items: center;
          gap: var(--spacing-xs);
        }
        
        .detail-item input {
          width: 60px;
          padding: 4px;
          text-align: center;
        }
        
        .exercise-actions {
          display: flex;
          gap: var(--spacing-xs);
        }
        
        .action-btn {
          padding: 5px;
          font-size: 0.8rem;
        }
        
        .exercise-filters {
          display: flex;
          gap: var(--spacing-md);
          margin-bottom: var(--spacing-md);
        }
        
        .search-input-container {
          position: relative;
          flex: 1;
        }
        
        .search-input-container .fa-filter {
          position: absolute;
          left: 10px;
          top: 50%;
          transform: translateY(-50%);
          color: var(--text-muted);
        }
        
        .search-input {
          padding-left: 30px;
          width: 100%;
        }
        
        .available-exercises {
          max-height: 500px;
          overflow-y: auto;
        }
        
        .exercise-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: var(--spacing-md);
          border-bottom: 1px solid var(--border-color);
        }
        
        .exercise-info {
          flex: 1;
        }
        
        .exercise-info h4 {
          margin: 0;
          margin-bottom: var(--spacing-xs);
        }
        
        .exercise-metadata {
          font-size: 0.85rem;
          color: var(--text-muted);
          display: flex;
          gap: var(--spacing-md);
          margin-bottom: var(--spacing-xs);
        }
        
        .exercise-description {
          font-size: 0.9rem;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
          text-overflow: ellipsis;
        }
        
        .empty-state {
          padding: var(--spacing-lg);
          text-align: center;
          color: var(--text-muted);
        }
        
        .hint {
          font-size: 0.9rem;
          margin-top: var(--spacing-md);
        }
        
        .mt-md {
          margin-top: var(--spacing-md);
        }
        
        .mt-lg {
          margin-top: var(--spacing-lg);
        }
        
        .text-center {
          text-align: center;
        }
        
        .loading {
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
      `}</style>
    </div>
  );
};

export default NewRoutinePage;
