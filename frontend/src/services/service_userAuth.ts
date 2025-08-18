import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import toast from 'react-hot-toast';
import * as io from 'socket.io-client';
import { SOCKET_BASE_URL} from '../configs/contants';
import {axiosInstance} from '../lib/axios';
import axios from 'axios';

import type { AuthContextType , AuthUser, SocketError } from '../configs/types';
import type { SignUpData, SignInData, UpdateProfileData } from '../configs/types';


export const userAuthService = create<AuthContextType>()(persist(
    (set, get) => ({
    authUser: null,
    isSigningUp: false,
    isLoggingIn: false,
    isUpdatingProfile: false,
    isCheckingAuth: true,
    onlineUsers: [] as string[],
    socket: null,
    token: null,

    checkAuth: async () => {
        try {
            const res  = await axiosInstance.get<AuthUser>("/auth/check-auth");
            set({authUser: res.data});
            get().connectSocket();
        } catch (error) {
            console.log("Error in Auth Check", error);
            // set({authUser: null});
            set({ authUser: null, token: null });
        } finally {
            set({isCheckingAuth: false});
        }
    },

    signUp: async (signUpData: SignUpData) => {
        try {
            set({isSigningUp: true});
            const res = await axiosInstance.post("/signup", signUpData);            
            if (res.data.signupStatus || res.data.message === "User created successfully") {
                // set({authUser: res.data.user});
                toast.success("Account created successfully");
                // get().connectSocket();
            }
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
    
            // set({authUser: user});
            // toast.success("Login successful");
            // get().connectSocket();
            if (user?._id) {
                // console.log("Have running this herre..... ", user._id);
                set({authUser: user, token: token});
                toast.success("Login successful");
                get().connectSocket();
            }
            console.log("Connecting socket with userId:", user?._id);
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
            set({authUser: null, token: null});
            toast.success("Logout successful");
            get().disconnectSocket();
            // Set in localStorage to null
            localStorage.removeItem("token");
            localStorage.removeItem("authUser");
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
        // finally {
        //     set({ authUser: null, token: null });
        //     localStorage.removeItem("token");
        //     get().disconnectSocket();
        //     toast.success("Logout successful");
        // }
    },

    updateProfile: async (updateProfileData: UpdateProfileData) => {
        try {
            const userId = get().authUser?._id;
            const res = await axiosInstance.put(`/users/${userId}/update-profile`, updateProfileData);
            set({authUser: res.data});
            set({isUpdatingProfile: true});
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
        // Version protocol 3, 4
    // connectSocket: () => {
    //     const {authUser} = get();
    //     console.log("Check socket:", authUser);
    //     console.log("Check socket:", authUser?._id);
    //     console.log("Check socket:", get().socket?.connected);
    //     if (!authUser || !authUser._id || get().socket?.connected) return;
    //     console.log("Passing check 1......")
    //     const socket = io(`${SOCKET_BASE_URL}?userId=${authUser._id}`, {
    //         transports: ["websocket"],
    //         autoConnect: true,
    //         // path: "/socket.io",
    //         // query: {
    //         //     userId: authUser._id,
    //         // },
    //         // auth: {
    //         //     token: localStorage.getItem("token") || "",
    //         // },
    //         // extraHeaders: {
    //         //     "X-Socketio-Auth": JSON.stringify({ token: localStorage.getItem("token") }),
    //         // },
    //         withCredentials: true,
    //     });

    //     console.log("Passing check 3......")
    //     // socket.connect();
    //     console.log("Passing check 4......")

    //     set({socket: socket });
    //     socket.on("connect", () => {
    //         console.log("✅ Connected with id:", socket.id);
    //     });
    //     socket.on("connect_error", (err) => {
    //         console.error("❌ connect_error:", err.message, err.cause, err.stack);
    //     });
    //     socket.on("disconnect", (reason) => {
    //         console.warn("⚠️ Socket disconnected:", reason);
    //     });
    //     socket.on("getOnlineUsers", (userIds: string[]) => {
    //         console.log("Online users Tessting in Socket:", userIds);
    //         set({onlineUsers: userIds});
    //     });

    //     console.log("Passing check 5......")
    //     console.log("Socket connected!", get().onlineUsers);

    // },
    connectSocket: () => {
        const { authUser } = get();
        if (!authUser || !authUser._id || get().socket?.connected) return;

        const socket = io.connect(`${SOCKET_BASE_URL}/?userId=${authUser._id}`, {
            // query: { userId: authUser._id },
            transports: ["websocket"],
            path: "/socket.io"
        });

        set({ socket: socket });

        socket.on("connect", () => {
            console.log("✅ Connected with id:", socket.id);
        });

        // socket.on("connect_error", (err: SocketError) => {
        //     console.error("❌ connect_error:", err.message);
        // });
        socket.addEventListener("getOnlineUsers", (event: MessageEvent) => {
            console.log("[getOnlineUsers]: Online users:", event.data);
            set({ onlineUsers: event.data });
        });
        // Kiểm tra tất cả event
        socket.on("socket.io/getOnlineUsers", (onlineUsers: MessageEvent) => {
            console.log("[getOnlineUsers]: Online users:", onlineUsers);
            // set({ onlineUsers: onlineUsers });
        });
        socket.on("disconnect", (reason: SocketError) => {
            console.warn("⚠️ Socket disconnected:", reason);
        });

    },

    disconnectSocket: () => {
        console.log("Socket disconnected!", get().socket);
        const socket = get().socket;
        if (socket && socket.connected) {
            socket.disconnect();
            set({socket: null});
        }
    },

}),
    {
      name: "auth-storage", // key trong localStorage
      partialize: (state) => ({ authUser: state.authUser, token: state.token }), 
    }
));

