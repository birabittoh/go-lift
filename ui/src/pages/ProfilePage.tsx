import { useState, useEffect } from 'react';
import { useAppContext } from '../context/AppContext';
import type { User } from '../types/models';
import { FaUser, FaSave, FaTimes } from 'react-icons/fa';

function parseBirthDate(dateString: string) {
  return new Date(dateString).toISOString().split('T')[0]; // Format to YYYY-MM-DD
}

const ProfilePage = () => {
  const { user, updateUser, isLoading, error: contextError } = useAppContext();
  const [formData, setFormData] = useState<User | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  useEffect(() => {
    if (user && !formData) {
      setFormData({ ...user });
    }
  }, [user, formData]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    
    setFormData(prev => {
      if (!prev) return prev;
      return { ...prev, [name]: value };
    });
  };

  const handleEditToggle = () => {
    if (isEditing) {
      // Cancel edit - revert changes
      setFormData(user ? { ...user } : null);
    }
    setIsEditing(!isEditing);
    setError(null);
    setSuccessMessage(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData) return;
    
    try {
      setError(null);
      await updateUser(formData);
      setSuccessMessage('Profile updated successfully!');
      setIsEditing(false);
      
      // Clear success message after a few seconds
      setTimeout(() => {
        setSuccessMessage(null);
      }, 3000);
    } catch {
      setError('Failed to update profile. Please try again.');
    }
  };

  if (isLoading) {
    return (
      <div className="page profile-page">
        <h1>Profile</h1>
        <div className="loading">Loading profile data...</div>
      </div>
    );
  }

  if (!formData) {
    return (
      <div className="page profile-page">
        <h1>Profile</h1>
        <div className="error-message">Could not load profile data.</div>
      </div>
    );
  }

  return (
    <div className="page profile-page">
      <div className="profile-header">
        <h1>Your Profile</h1>
        
        <button 
          onClick={handleEditToggle}
          className={`btn ${isEditing ? 'btn-danger' : 'btn-secondary'}`}
        >
          {isEditing ? (
            <>
              <FaTimes /> Cancel
            </>
          ) : (
            <>Edit Profile</>
          )}
        </button>
      </div>

      {contextError && <div className="error-message">{contextError}</div>}
      {error && <div className="error-message">{error}</div>}
      {successMessage && <div className="success-message">{successMessage}</div>}

      <div className="card">
        <form onSubmit={handleSubmit}>
          <div className="profile-avatar">
            <div className="avatar-circle">
              <FaUser size={40} />
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="name">Name</label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
              disabled={!isEditing}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="isFemale">Gender</label>
            <select
              id="isFemale"
              name="isFemale"
              value={formData.isFemale.toString()}
              onChange={handleInputChange}
              disabled={!isEditing}
              required
            >
              <option value="false">Male</option>
              <option value="true">Female</option>
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="weight">Weight (kg)</label>
            <input
              type="number"
              id="weight"
              name="weight"
              min="20"
              max="300"
              step="1"
              value={formData.weight ?? 0}
              onChange={handleInputChange}
              disabled={!isEditing}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="height">Height (cm)</label>
            <input
              type="number"
              id="height"
              name="height"
              min="0"
              max="300"
              step="1"
              value={formData.height ?? 0}
              onChange={handleInputChange}
              disabled={!isEditing}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="birthDate">Date of Birth</label>
            <input
              type="date"
              id="birthDate"
              name="birthDate"
              value={parseBirthDate(formData.birthDate ?? '')}
              onChange={handleInputChange}
              disabled={!isEditing}
              required
            />
          </div>

          {isEditing && (
            <div className="form-actions">
              <button type="submit" className="btn btn-primary btn-block">
                <FaSave /> Save Changes
              </button>
            </div>
          )}
        </form>
      </div>
      
      <style>{`
        .profile-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: var(--spacing-lg);
        }
        
        .profile-avatar {
          display: flex;
          justify-content: center;
          margin-bottom: var(--spacing-xl);
        }
        
        .avatar-circle {
          width: 100px;
          height: 100px;
          border-radius: 50%;
          background-color: var(--light-gray);
          display: flex;
          align-items: center;
          justify-content: center;
          color: var(--dark-gray);
        }
        
        .form-actions {
          margin-top: var(--spacing-lg);
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

export default ProfilePage;
