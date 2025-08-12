import { useEffect } from "react";

import { Route, Routes, Navigate } from "react-router-dom";
import { userAuthService } from "./services/service_userAuth";
import { useThemeChange } from "./services/service_changeTheme";

import Navbar  from "./components/navbar"
import HomePage from "./pages/homePage";
import SignUpPage from "./pages/signupPage";
import LoginPage from "./pages/loginPage";
import SettingsPage from "./pages/settingsPage";
import ProfilePage from "./pages/profilePage";


import { Loader } from "lucide-react"; // Icon for loading state
import { Toaster } from "react-hot-toast"; // For displaying toast notifications


const App = () => {
    const { authUser, checkAuth, isCheckingAuth, onlineUsers } = userAuthService((state) => state);
    const { theme } = useThemeChange((state) => state);

    console.log("Online users: ", { onlineUsers });
    useEffect(() => {
        checkAuth();
    }, [checkAuth]);

    console.log("userAuth state: ", { authUser, isCheckingAuth });

    if (isCheckingAuth && !authUser) {
        return (
            <div className="flex items-center justify-center h-screen">
                <Loader className="size-10 animate-spin" />
            </div>
        )
    }

    return (
        <div data-theme={theme} className="w-screen">
            <Navbar />
            <Routes>
                <Route path="/" element={authUser ? <HomePage /> : <Navigate to="/login" />} />
                <Route path="/signup" element={authUser ? <Navigate to="/" /> : <SignUpPage />} />
                <Route path="/login" element={authUser ? <Navigate to="/" /> : <LoginPage />} />
                <Route path="/settings" element={authUser ? <SettingsPage /> : <Navigate to="/login" />} />
                <Route path="/profile" element={authUser ? <ProfilePage /> : <Navigate to="/login" />} />
            </Routes>
            <Toaster />
        </div>
    );
};

export default App;