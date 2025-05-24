import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import type { User } from '../types/models';
import { userService } from '../services/api';

interface AppContextType {
  user: User | null;
  isLoading: boolean;
  error: string | null;
  updateUser: (profile: User) => Promise<void>;
}

const defaultProfile: User = {
  name: 'User',
  isFemale: false,
  weight: 85,
  height: 180,
  birthDate: "01/01/1990",
};

const AppContext = createContext<AppContextType>({
  user: defaultProfile,
  isLoading: false,
  error: null,
  updateUser: async () => {},
});

export const useAppContext = () => useContext(AppContext);

interface AppProviderProps {
  children: ReactNode;
}

export const AppProvider = ({ children }: AppProviderProps) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadProfile = async () => {
      try {
        setIsLoading(true);
        const profile = await userService.get(1);
        setUser(profile);
        setError(null);
      } catch (err) {
        console.error('Failed to load user profile:', err);
        // For the single-user mode, create a default profile if none exists
        setUser(defaultProfile);
        setError('Could not load profile. Using default settings.');
      } finally {
        setIsLoading(false);
      }
    };

    loadProfile();
  }, []);

  const updateUser = async (profile: User) => {
    try {
      setIsLoading(true);
      const updatedProfile = await userService.update(1, profile);
      setUser(updatedProfile);
      setError(null);
    } catch (err) {
      console.error('Failed to update profile:', err);
      setError('Failed to update profile. Please try again.');
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AppContext.Provider value={{ user, isLoading, error, updateUser }}>
      {children}
    </AppContext.Provider>
  );
};
