package ws

import (
	"log"
	"net/url"

	socketio "github.com/googollee/go-socket.io"
)

// userSocketMap stores the mapping between user ID and socket ID
var userSocketMap = make(map[string]string)

func GetReceiverSocketId(userId string) string {
	return userSocketMap[userId]
}

func NewServer() *socketio.Server {
	server := socketio.NewServer(nil)

	// Handle new connections
	server.OnConnect("/", func(s socketio.Conn) error {

		log.Println("A user connected:", s.ID())

		//Get userId from the handshake query
		urlValues, err := url.ParseQuery(s.URL().RawQuery)
		if err != nil {
			log.Printf("Error parsing URL query from handshake for socket %s: %v", s.ID(), err)
			return err
		}
		userId := urlValues.Get("userId")

		if userId != "" {
			userSocketMap[userId] = s.ID()
			log.Println("User mapped: ", userId, "->", s.ID())
		}
		//Broadcast the updated list of online users to all clients.
		onlineUsers := make([]string, 0, len(userSocketMap))
		for k := range userSocketMap {
			onlineUsers = append(onlineUsers, k)
		}
		server.BroadcastToNamespace("/", "getOnlineUsers", onlineUsers)
		return nil

	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("Socket error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("A user disconnected:", s.ID(), "Reason:", reason)

		// Find the UserId associated with the disconnected socket.
		var userIdToRemove string
		for id, socketId := range userSocketMap {
			if socketId == s.ID() {
				userIdToRemove = id
				break
			}
		}
		// If a userId was found, remove it from the map
		if userIdToRemove != "" {
			delete(userSocketMap, userIdToRemove)
			log.Println("User unmapped: ", userIdToRemove, "->", s.ID())
		}

		// Broadcast the updated list of online users to all clients.
		onlineUsers := make([]string, 0, len(userSocketMap))
		for k := range userSocketMap {
			onlineUsers = append(onlineUsers, k)
		}

		server.BroadcastToNamespace("/", "getOnlineUsers", onlineUsers)
	})

	return server
}

// ErrServerClosed
type ErrServerClosed struct{}

func (e *ErrServerClosed) Error() string {
	return "server closed"
}
