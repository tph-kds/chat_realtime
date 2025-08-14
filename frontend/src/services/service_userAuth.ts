import { create } from 'zustand';
import toast from 'react-hot-toast';
import {io} from 'socket.io-client';
import { SOCKET_BASE_URL} from '../configs/contants';
import {axiosInstance} from '../lib/axios';
import axios from 'axios';

import type { AuthContextType } from '../configs/types';
import type { SignUpData, SignInData, UpdateProfileData } from '../configs/types';


export const userAuthService = create<AuthContextType>((set, get) => ({
    authUser: null,
    isSigningUp: false,
    isLoggingIn: false,
    isUpdatingProfile: false,
    isCheckingAuth: true,
    onlineUsers: [],
    socket: null,

    checkAuth: async () => {
        try {
            const res  = await axiosInstance.get<AuthContextType>("/auth/check-auth");
            set({authUser: res.data.authUser});
            get().connectSocket();
        } catch (error) {
            console.log("Error in Auth Check", error);
            set({authUser: null});
        } finally {
            set({isCheckingAuth: false});
        }
    },

    signUp: async (signUpData: SignUpData) => {
        try {
            set({isSigningUp: true});
            const res = await axiosInstance.post("/signup", signUpData);
            set({authUser: res.data});
            toast.success("Account created successfully");
            get().connectSocket();
            // return res;
        } catch (error) {
            toast.error("Account creation failed");

            if (axios.isAxiosError(error)) {
                toast.error(error.response?.data?.message || "Server error occurred");
            } else if (error instanceof Error) {
                toast.error(error.message);
            } else {
                toast.error("An unknown error occurred");
            }

            console.log("Error in signup", error);
            // return error;
        } finally {
            set({isSigningUp: false});
        }
    },

    signIn: async (signInData: SignInData) => {
        try {
            set({isLoggingIn: true});
            const res = await axiosInstance.post("/login", signInData);
            // Get Token  and save it
            const { token , user } = res.data;
            // Lưu token vào localStorage
            localStorage.setItem("token", token);
    
            set({authUser: user});
            toast.success("Login successful");
            get().connectSocket();
            // return res;
        } catch (error) {
            toast.error("Login failed");

            if (axios.isAxiosError(error)) {
                toast.error(error.response?.data?.message || "Server error occurred");
            } else if (error instanceof Error) {
                toast.error(error.message);
            } else {
                toast.error("An unknown error occurred");
            }

            console.log("Error in signin", error);
            // return error;
        } finally {
            set({isLoggingIn: false});
        }
    },

    logOut: async () => {
        try {
            await axiosInstance.post("/logout");
            set({authUser: null});
            toast.success("Logout successful");
            get().disconnectSocket();
            // return res;
        } catch (error) {
            toast.error("Logout failed");

            if (axios.isAxiosError(error)) {
                toast.error(error.response?.data?.message || "Server error occurred");
            } else if (error instanceof Error) {
                toast.error(error.message);
            } else {
                toast.error("An unknown error occurred");
            }

            console.log("Error in logout", error);
        }
    },

    updateProfile: async (updateProfileData: UpdateProfileData) => {
        set({isUpdatingProfile: true});
        try {
            const userId = get().authUser?._id;
            const res = await axiosInstance.put(`/users/${userId}/update-profile`, updateProfileData);
            set({authUser: res.data});
            toast.success("Profile updated successfully");
            // return res;
        } catch (error) {
            toast.error("Profile update failed");

            if (axios.isAxiosError(error)) {
                toast.error(error.response?.data?.message || "Server error occurred");
            } else if (error instanceof Error) {
                toast.error(error.message);
            } else {
                toast.error("An unknown error occurred");
            }

            console.log("Error in update profile", error);
        } finally {
            set({isUpdatingProfile: false});

        }
    },

    connectSocket: () => {
        const {authUser} = get();
        if (!authUser || get().socket?.connected) return;

        const socket = io(SOCKET_BASE_URL, {
            // transports: ["websocket"],
            // autoConnect: true,
            // path: "/socket.io",
            query: {
                userId: authUser._id,
            },
        });

        socket.connect();

        set({socket: socket });
        socket.on("getOnlineUsers", (userIds) => {
            set({onlineUsers: userIds});
        });
    },

    disconnectSocket: () => {
        const socket = get().socket;
        if (socket && socket.connected) {
            socket.disconnect();
            set({socket: null});
        }
    },

}));

