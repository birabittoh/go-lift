import { useEffect, useState } from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import HomePage from "./pages/HomePage";
import WorkoutsPage from "./pages/WorkoutsPage";
import NewRoutinePage from "./pages/NewRoutinePage";
import NewExercisePage from "./pages/NewExercisePage";
import NewWorkoutPage from "./pages/NewWorkoutPage";
import ProfilePage from "./pages/ProfilePage";
import BottomNav from "./components/BottomNav";
import { AppProvider } from "./context/AppContext";
import "./App.css";

function App() {
  const [isMobile, setIsMobile] = useState(window.innerWidth < 768);

  // Check for device width to determine if mobile or desktop
  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < 768);
    };

    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  return (
    <AppProvider>
      <Router>
        <div className={`app-container ${isMobile ? 'mobile-layout' : 'desktop-layout'}`}>
          <div className="content-area">
            <Routes>
              <Route path="/" element={<Navigate to="/home" replace />} />
              <Route path="/home" element={<HomePage />} />
              <Route path="/workouts" element={<WorkoutsPage />} />
              <Route path="/new-routine" element={<NewRoutinePage />} />
              <Route path="/new-exercise" element={<NewExercisePage />} />
              <Route path="/new-workout" element={<NewWorkoutPage />} />
              <Route path="/profile" element={<ProfilePage />} />
            </Routes>
          </div>
          <BottomNav isMobile={isMobile} />
        </div>
      </Router>
    </AppProvider>
  );
}

export default App;
