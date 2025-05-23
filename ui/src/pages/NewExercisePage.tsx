import { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { FaSave, FaArrowLeft, FaTrash, FaTimes } from 'react-icons/fa';
import type { Exercise, Equipment, MuscleGroup } from '../types/models';
import { ExerciseService, EquipmentService, MuscleGroupService } from '../services/api';

const NewExercisePage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  
  const [exerciseToEdit, setExerciseToEdit] = useState<Exercise | null>(null);
  const [isEditMode, setIsEditMode] = useState<boolean>(false);
  
  // Available equipment and muscle groups from API
  const [availableEquipment, setAvailableEquipment] = useState<Equipment[]>([]);
  const [availableMuscleGroups, setAvailableMuscleGroups] = useState<MuscleGroup[]>([]);
  
  // Selected items
  const [selectedEquipment, setSelectedEquipment] = useState<Equipment[]>([]);
  const [selectedMuscleGroups, setSelectedMuscleGroups] = useState<MuscleGroup[]>([]);
  
  // Form state
  const [formData, setFormData] = useState<Exercise>({
    name: '',
    description: '',
    equipment: [],
    muscleGroups: [],
    sets: [],
  });
  
  // Fetch equipment and muscle groups from API
  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      try {
        const [equipmentData, muscleGroupsData] = await Promise.all([
          EquipmentService.getAll(),
          MuscleGroupService.getAll()
        ]);
        setAvailableEquipment(equipmentData);
        setAvailableMuscleGroups(muscleGroupsData);
      } catch (err) {
        console.error('Failed to fetch data:', err);
        setError('Failed to load equipment and muscle groups');
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchData();
  }, []);
  
  // Check if we're editing an existing exercise
  useEffect(() => {
    if (location.state && location.state.editExercise) {
      const exercise = location.state.editExercise as Exercise;
      setFormData(exercise);
      setSelectedEquipment(exercise.equipment || []);
      setSelectedMuscleGroups(exercise.muscleGroups || []);
      setExerciseToEdit(exercise);
      setIsEditMode(true);
    }
  }, [location]);
  
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  };
  
  const handleEquipmentSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const equipmentId = parseInt(e.target.value);
    const equipment = availableEquipment.find(eq => eq.id === equipmentId);
    
    if (equipment && !selectedEquipment.some(e => e.id === equipment.id)) {
      const updatedEquipment = [...selectedEquipment, equipment];
      setSelectedEquipment(updatedEquipment);
      setFormData(prev => ({
        ...prev,
        equipment: updatedEquipment
      }));
    }
    
    // Reset select to default
    e.target.value = '';
  };
  
  const removeEquipment = (equipmentId?: number) => {
    if (!equipmentId) return;
    const updatedEquipment = selectedEquipment.filter(e => e.id !== equipmentId);
    setSelectedEquipment(updatedEquipment);
    setFormData(prev => ({
      ...prev,
      equipment: updatedEquipment
    }));
  };
  
  const handleMuscleGroupSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const muscleGroupId = parseInt(e.target.value);
    const muscleGroup = availableMuscleGroups.find(mg => mg.id === muscleGroupId);
    
    if (muscleGroup && !selectedMuscleGroups.some(mg => mg.id === muscleGroup.id)) {
      const updatedMuscleGroups = [...selectedMuscleGroups, muscleGroup];
      setSelectedMuscleGroups(updatedMuscleGroups);
      setFormData(prev => ({
        ...prev,
        muscleGroups: updatedMuscleGroups
      }));
    }
    
    // Reset select to default
    e.target.value = '';
  };
  
  const removeMuscleGroup = (muscleGroupId?: number) => {
    if (!muscleGroupId) return;
    const updatedMuscleGroups = selectedMuscleGroups.filter(mg => mg.id !== muscleGroupId);
    setSelectedMuscleGroups(updatedMuscleGroups);
    setFormData(prev => ({
      ...prev,
      muscleGroups: updatedMuscleGroups
    }));
  };
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    
    try {
      if (formData.muscleGroups.length === 0) {
        throw new Error('Please select at least one muscle group');
      }
      
      if (isEditMode && exerciseToEdit?.id) {
        // Update existing exercise
        await ExerciseService.update(exerciseToEdit.id, formData);
        setSuccessMessage('Exercise updated successfully!');
      } else {
        // Create new exercise
        await ExerciseService.create(formData);
        setSuccessMessage('Exercise created successfully!');
        
        // Reset form if creating new
        if (!isEditMode) {
          setFormData({
            name: '',
            description: '',
            equipment: [],
            muscleGroups: [],
            sets: [],
          });
          setSelectedEquipment([]);
          setSelectedMuscleGroups([]);
        }
      }
      
      // Show success message briefly then redirect
      setTimeout(() => {
        navigate('/workouts');
      }, 1500);
    } catch (err: unknown) {
      console.error('Failed to save exercise:', err);
      setError('Failed to save exercise. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };
  
  const handleDelete = async () => {
    if (!exerciseToEdit?.id || !confirm('Are you sure you want to delete this exercise?')) {
      return;
    }
    
    setIsLoading(true);
    try {
      await ExerciseService.delete(exerciseToEdit.id);
      setSuccessMessage('Exercise deleted successfully!');
      
      // Redirect after deletion
      setTimeout(() => {
        navigate('/workouts');
      }, 1500);
    } catch (err) {
      console.error('Failed to delete exercise:', err);
      setError('Failed to delete exercise. It might be used in one or more routines.');
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <div className="page new-exercise-page">
      <div className="page-header">
        <button 
          onClick={() => navigate(-1)} 
          className="btn btn-secondary back-button"
        >
          <FaArrowLeft /> Back
        </button>
        <h1>{isEditMode ? 'Edit Exercise' : 'New Exercise'}</h1>
      </div>
      
      {error && <div className="error-message">{error}</div>}
      {successMessage && <div className="success-message">{successMessage}</div>}
      
      <div className="card">
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="name">Exercise Name*</label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
              placeholder="e.g. Bench Press"
              required
              disabled={isLoading}
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              name="description"
              value={formData.description}
              onChange={handleInputChange}
              placeholder="Describe how to perform this exercise properly..."
              rows={4}
              disabled={isLoading}
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="muscleGroups">Muscle Groups*</label>
            <select
              id="muscleGroups"
              name="muscleGroups"
              onChange={handleMuscleGroupSelect}
              disabled={isLoading}
            >
              <option value="">Select muscle group...</option>
              {availableMuscleGroups.map(group => (
                <option key={group.id} value={group.id}>{group.name}</option>
              ))}
            </select>
            
            <div className="tags-container">
              {selectedMuscleGroups.map(group => (
                <div key={group.id} className="tag">
                  {group.name}
                  <button 
                    type="button" 
                    className="tag-remove" 
                    onClick={() => removeMuscleGroup(group.id)}
                    disabled={isLoading}
                  >
                    <FaTimes />
                  </button>
                </div>
              ))}
            </div>
          </div>
          
          <div className="form-group">
            <label htmlFor="equipment">Equipment</label>
            <select
              id="equipment"
              name="equipment"
              onChange={handleEquipmentSelect}
              disabled={isLoading}
            >
              <option value="">Select equipment...</option>
              {availableEquipment.map(eq => (
                <option key={eq.id} value={eq.id}>{eq.name}</option>
              ))}
            </select>
            
            <div className="tags-container">
              {selectedEquipment.map(eq => (
                <div key={eq.id} className="tag">
                  {eq.name}
                  <button 
                    type="button" 
                    className="tag-remove" 
                    onClick={() => removeEquipment(eq.id)}
                    disabled={isLoading}
                  >
                    <FaTimes />
                  </button>
                </div>
              ))}
            </div>
          </div>
          
          <div className="form-actions">
            <button 
              type="submit" 
              className="btn btn-primary btn-block"
              disabled={isLoading}
            >
              <FaSave /> {isEditMode ? 'Update Exercise' : 'Save Exercise'}
            </button>
            
            {isEditMode && exerciseToEdit?.id && (
              <button
                type="button"
                className="btn btn-danger btn-block mt-md"
                onClick={handleDelete}
                disabled={isLoading}
              >
                <FaTrash /> Delete Exercise
              </button>
            )}
          </div>
        </form>
      </div>
      
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
        
        .form-actions {
          margin-top: var(--spacing-xl);
        }
        
        textarea {
          resize: vertical;
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
        
        .tags-container {
          display: flex;
          flex-wrap: wrap;
          gap: var(--spacing-sm);
          margin-top: var(--spacing-sm);
        }
        
        .tag {
          display: flex;
          align-items: center;
          background-color: var(--primary-color-light);
          color: var(--primary-color-dark);
          padding: var(--spacing-xs) var(--spacing-sm);
          border-radius: var(--border-radius);
          font-size: 0.9rem;
        }
        
        .tag-remove {
          border: none;
          background: none;
          color: var(--primary-color-dark);
          margin-left: var(--spacing-xs);
          padding: 0;
          font-size: 0.8rem;
          display: flex;
          align-items: center;
          cursor: pointer;
        }
        
        .mt-md {
          margin-top: var(--spacing-md);
        }
      `}</style>
    </div>
  );
};

export default NewExercisePage;
