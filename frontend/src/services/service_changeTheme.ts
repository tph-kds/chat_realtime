import { create } from 'zustand';
import { toast } from 'react-hot-toast';
import type { Theme, ThemeState, ThemeActions } from '../configs/types';


/**
 * Safely retrieves the theme from localStorage.
 * Falls back to a default theme if no valid theme is found.
 * @returns {Theme} The stored or default theme.
 */

const getInitialTheme = () : Theme => {
    try {
        const storedTheme = localStorage.getItem("chat-theme");
        if (storedTheme && ["light", "dark", "coffee", "corporate"].includes(storedTheme)) {
            return storedTheme as Theme;
        } else {
            return "light";
        }
    } catch (error) {
        console.error("Error retrieving theme from localStorage:", error);
        // return "light";
    } 
    return "coffee";
}

export const useThemeChange = create<ThemeState & ThemeActions>((set) => ({
    theme: getInitialTheme(),
    setTheme: (theme: Theme) => {
        try {
            localStorage.setItem("chat-theme", theme);
        } catch (error) {
            console.error("Error setting theme in localStorage:", error);
            toast.error("Error setting theme in localStorage");
        }
        set({ theme });
    },
}));