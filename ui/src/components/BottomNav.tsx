import { NavLink } from "react-router-dom";
import { 
  FaHome, 
  FaDumbbell, 
  FaUserAlt, 
  FaPlus
} from "react-icons/fa";
import "./BottomNav.css";

interface BottomNavProps {
  isMobile: boolean;
}

const BottomNav = ({ isMobile }: BottomNavProps) => {
  if (isMobile) {
    return (
      <nav className="bottom-nav mobile">
        <NavLink to="/home" className={({ isActive }) => isActive ? "active" : ""}>
          <FaHome size={20} />
          <span>Home</span>
        </NavLink>
        <NavLink to="/workouts" className={({ isActive }) => isActive ? "active" : ""}>
          <FaDumbbell size={20} />
          <span>Workouts</span>
        </NavLink>
        <NavLink to="/profile" className={({ isActive }) => isActive ? "active" : ""}>
          <FaUserAlt size={20} />
          <span>Profile</span>
        </NavLink>
      </nav>
    );
  }
  
  return (
    <nav className="side-nav desktop">
      <div className="app-title">Go Lift</div>
      <div className="nav-links">
        <NavLink to="/home" className={({ isActive }) => isActive ? "active" : ""}>
          <FaHome size={20} />
          <span>Home</span>
        </NavLink>
        <NavLink to="/workouts" className={({ isActive }) => isActive ? "active" : ""}>
          <FaDumbbell size={20} />
          <span>Workouts</span>
        </NavLink>
        <div className="sub-menu">
          <div className="sub-menu-title">Create New</div>
          <NavLink to="/new-exercise" className={({ isActive }) => isActive ? "active" : ""}>
            <FaPlus size={16} />
            <span>New Exercise</span>
          </NavLink>
          <NavLink to="/new-routine" className={({ isActive }) => isActive ? "active" : ""}>
            <FaPlus size={16} />
            <span>New Routine</span>
          </NavLink>
          <NavLink to="/new-workout" className={({ isActive }) => isActive ? "active" : ""}>
            <FaPlus size={16} />
            <span>New Workout</span>
          </NavLink>
        </div>
        <NavLink to="/profile" className={({ isActive }) => isActive ? "active" : ""}>
          <FaUserAlt size={20} />
          <span>Profile</span>
        </NavLink>
      </div>
    </nav>
  );
};

export default BottomNav;
