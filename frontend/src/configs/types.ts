
import { Socket } from "socket.io-client";

export interface BaseAttributes {
    _id: string;
    createdAt: Date;
    updatedAt: Date;
}

export interface User extends BaseAttributes {
//   _id: string;
  first_name: string;
  last_name: string;
  email: string;
  avatar: string;
//   createdAt: Date;
//   updatedAt: Date;
  token: string;
  isOnline: boolean;
}

export interface Message extends BaseAttributes {
    // _id: string;
    chatId: string;
    senderId: string;
    text: string;
    // createdAt: Date;
    // updatedAt: Date;
}

export interface Chat extends BaseAttributes {
    // _id: string;
    members: string[];
    // createdAt: Date;
    // updatedAt: Date;
}



// ######################## THE AUTHENTICATED USER OBJECT TYPE DECLARATION #######################
export interface AuthUser {
    _id: string;
    userName: string;
    fullName: string;
    email: string;
    profilePic: string;
    createdAt: string;
    updatedAt: string;
}

export interface SignUpData extends Record<string, unknown> {
    email: string;
    password: string;
}

export interface SignInData extends Record<string, unknown> {
    // createdAt: Date;
    // updatedAt: Date;
    email: string;
    password: string;
}

export interface UpdateProfileData extends Record<string, unknown> {
    base64Image: string;
    createdAt: string;
    updatedAt: string;
}

export interface APIError {
    message: string;
}


export interface AuthState {
    authUser: AuthUser | null;
    isSigningUp: boolean | false;
    isLoggingIn: boolean | false;
    isUpdatingProfile: boolean | false;
    isCheckingAuth: boolean | true;
    onlineUsers: string[];
    socket: Socket |null;
}

export interface AuthActions {
    checkAuth: () => Promise<void>;
    signUp: (signUpData: SignUpData) => Promise<void>;
    signIn: (signInData: SignInData) => Promise<void>;
    logOut: () => Promise<void>;
    updateProfile: (updateProfileData: UpdateProfileData) => Promise<void>;
    connectSocket: () => void;
    disconnectSocket: () => void;
}

type AuthContextType = AuthState & AuthActions;
export type { AuthContextType };



// ######################## THE CHAT USER OBJECT TYPE DECLARATION #######################
export interface ChatUser {
    _id: string;
    userName: string;
    fullName: string;
    profilePic: string;
}

export interface MessageChatUser {
    _id: string;
    receiverId: string;
    senderId: string;
    body: string;
    createdAt: Date;
    updatedAt: Date;
    image: string;
    text: string;
}

export interface SendMessageData {
    // receiverId: string;
    // body: string;
    image: string;
    text: string;
}

export interface ChatAPIError {
    message: string;
}

export interface ChatState {
    users: ChatUser[];
    messages: MessageChatUser[];
    selectedUser: ChatUser | null;
    isUsersLoading: boolean;
    isMessagesLoading: boolean;
}

export interface ChatActions {
    getUsers: () => Promise<void>;
    getMessages: (userId: string) => Promise<void>;
    sendMessage: (messageData: SendMessageData) => Promise<void>;
    subscribeToMessages: () => void;
    unsubscribeFromMessages: () => void;
    setSelectedUser: (user: ChatUser | null) => void;
    clearChatState: () => void;
}

type ChatContextType = ChatState & ChatActions;
export type { ChatContextType };


// ######################## THE THEME CHANGE OBJECT TYPE DECLARATION #######################
export type Theme = "light" | "dark" | "coffee" | "corporate";

export interface ThemeState {
    theme: Theme;
}

export interface ThemeActions {
    setTheme: (theme: Theme) => void;
}

type ThemeContextType = ThemeState & ThemeActions;
export type { ThemeContextType };

// ######################## THE LOGINPAGE OBJECT TYPE DECLARATION #######################
export interface LoginPageData extends Record<string, unknown> {
    email: string;
    password: string;
}

type LogInPageDataContextType = LoginPageData;
export type { LogInPageDataContextType };


// ######################## THE SIGNUP_PAGE OBJECT TYPE DECLARATION #######################
export interface SignUpPageData extends Record<string, unknown> {
    // full_name: string;
    first_name: string;
    last_name: string;
    email: string;
    password: string;
}

type SignUpPageDataContextType = SignUpPageData;
export type { SignUpPageDataContextType };



// ######################## THE COMPONENT OBJECTS TYPE DECLARATION #######################
export interface AuthImagePatternProps {
    title: string;
    subtitle: string;
}
// ######################## THE FORMAT MESSAGE OBJECTS TYPE DECLARATION #######################
export interface FormatMessageTimeProps {
    date: string;
}
