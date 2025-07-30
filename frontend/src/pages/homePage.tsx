import { userChatService } from "../services/service_userChat";
import SideBar from "../components/sideBar";
import NoChatSelected from "../components/noChatSelected";
import ChatContainer from "../components/chatContainer";



const HomePage: React.FC = () => {
    const { selectedUser } = userChatService();
    
    return (
        <div className="h-screen bg-base-200">
            <div className="flex items-center justify-center pt-20 px-4">
                <div className="bg-base-100 rounded-lg shadow-cl w-full max-w-6xl h-[calc(100vh-8rem)]">
                <div className="flex h-full rounded-lg overflow-hidden">
                    <SideBar />

                    {!selectedUser ? <NoChatSelected /> : <ChatContainer />}
                </div>
                </div>
            </div>
        </div>
    );

};

export default HomePage;