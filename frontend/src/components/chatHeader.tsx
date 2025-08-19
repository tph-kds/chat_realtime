import { X } from "lucide-react";
import { userAuthService } from "../services/service_userAuth";
import { userChatService } from "../services/service_userChat";
import { useEffect } from "react";


const ChatHeader = () => {
  const { selectedUser, setSelectedUser } = userChatService((state) => state);
  const onlineUsers  = userAuthService((state) => state.onlineUsers);
  // console.log("Online users: ", { onlineUsers });

      // Sử dụng useEffect để lắng nghe sự thay đổi
  useEffect(() => {
      console.log("Online users have been updated:", onlineUsers);
  }, [onlineUsers]);

  return (
    <div className="p-2.5 border-b border-base-300">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          {/* Avatar */}
          <div className="avatar">
            <div className="size-10 rounded-full relative">
              <img src={selectedUser?.profile_pic || "/avatar.png"} alt={selectedUser?.last_name} />
            </div>
          </div>

          {/* User info */}
          <div>
            <h3 className="font-medium">{selectedUser?.last_name || "User" }</h3>
            <p className="text-sm text-base-content/70">
              {onlineUsers.includes(selectedUser?._id as string) ? "Online" : "Offline"}
            </p>
          </div>
        </div>

        {/* Close button */}
        <button onClick={() => setSelectedUser(null)}>
          <X />
        </button>
      </div>
    </div>
  );
};
export default ChatHeader;