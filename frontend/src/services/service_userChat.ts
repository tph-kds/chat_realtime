import { create } from 'zustand';
import { AxiosError } from 'axios';
import toast from 'react-hot-toast';
import { axiosInstance } from '../lib/axios';
import { userAuthService } from './service_userAuth';
import type { ChatContextType, ChatState, ChatUser, MessageChatUser, SendMessageData, ChatAPIError } from '../configs/types';



const initialState: ChatState = {
    messages: [],
    users: [],
    selectedUser: null,
    isUsersLoading: false,
    isMessagesLoading: false,
};


export const userChatService = create<ChatContextType>((set, get) => ({
    ...initialState,
    
    /**
     * Fetches the list of users the authenticated user can chat with.
     */
    getUsers: async () => {
        set({isUsersLoading: true});
        try {
            const res = await axiosInstance.get<ChatUser[]>("/users");
            // console.log("No of users: ", res.data);
            set({users: res.data} );
            // return res;
        } catch (error) {
            const axiosError = error as AxiosError<ChatAPIError>;
            toast.error(axiosError.response?.data?.message || "Failed to fetch users.");
            console.error("Error fetching users:", error);
        } finally {
            set({isUsersLoading: false});
        }
    },


    /**
     * Fetches the message history with a specific user.
     * @param userId - The ID of the user to fetch messages for.
     */
    getMessages: async (userId: string) => {
        set({isMessagesLoading: true});
        try {
            const res = await axiosInstance.get<{messages: MessageChatUser[]}>(`/messages/${userId}`);
            set({messages: res.data.messages});
            // return res;
        } catch (error) {
            const axiosError = error as AxiosError<ChatAPIError>;
            toast.error(axiosError.response?.data?.message || "Failed to fetch messages.");
            console.error("Error fetching messages:", error);
        } finally {
            set({isMessagesLoading: false});
        }
    },

    /**
     * Sends a message to the currently selected user.
     * @param messageData - The content of the message to send.
     */
    sendMessage: async (messageData: SendMessageData) => {
        const { selectedUser , messages } = get();
        try {
            console.log("Sending message to: ", selectedUser);
            console.log("Sending message to: ", selectedUser?._id);

            const res = await axiosInstance.post(`/messages/send/${selectedUser?._id}`, messageData);
            console.log("Message sent: ", res.data.message);
            if (messages) {
                set({messages: [...messages, res.data.message]});
            } else {
                set({messages: [res.data.message]});
            }
            // set({messages: [...messages, res.data.message]});
            // return res;
        } catch (error) {
            const axiosError = error as AxiosError<ChatAPIError>;
            toast.error(axiosError.response?.data?.message || "Failed to send message.");
            console.error("Error sending message:", error);
        }
    },

        /**
     * Subscribes to real-time 'newMessage' events from the socket for the selected user.
     * Note: This should be called after a user is selected.
     */
    subscribeToMessages: () => {
        const socket = userAuthService.getState().socket;
        const { selectedUser } = get();

        if (!socket || !selectedUser) return;

        // Ensure we don't have duplicate listeners
        socket.off("newMessage");

        socket.on("newMessage", (message: MessageChatUser) => {
            const isFromSelectedUser = message.senderId === selectedUser._id;
            const isToCurrentUser = message.receiverId === userAuthService.getState().authUser?._id;
            if (isFromSelectedUser && isToCurrentUser) {
                set((state) => ({
                    messages: [...state.messages, message],
                }));
            }
        });
    },

    unsubscribeFromMessages: () => {
        const socket = userAuthService.getState().socket;
        if (socket) {
            socket.off("newMessage");
        }
    },

    /**
     * Sets the currently selected user for chatting.
     * @param user - The user object to select, or null to deselect.
     */
    setSelectedUser: (user) => {
        set({ selectedUser: user });
    },

    /**
     * Resets the chat store to its initial state. Useful on logout.
     */
    clearChatState: () => {
        get().unsubscribeFromMessages();
        set(initialState);
    }

}));