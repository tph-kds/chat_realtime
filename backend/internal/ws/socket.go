package ws

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var server *socketio.Server

// userSocketMap stores the mapping between user ID and socket ID
var userSocketMap = make(map[string]string)

func GetReceiverSocketId(userId string) string {
	return userSocketMap[userId]
}

func NewServerCustom() *socketio.Server {
	server := socketio.NewServer(
		&engineio.Options{
			Transports: []transport.Transport{
				&polling.Transport{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
				&websocket.Transport{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			},
			// RequestChecker: func(r *http.Request) (http.Header, error) {
			// 	userId := r.URL.Query().Get("userId")
			// 	if userId == "" {
			// 		return nil, fmt.Errorf("missing userId")
			// 	}
			// 	return true, nil
			// },
		})

	// Handle new connections
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Printf("[Socket] ✅ OnConnect triggered for socket ID: %s", s.ID())
		// return nil
		// 1. Lấy userId từ query string (vẫn như cũ)
		urlValues, err := url.ParseQuery(s.URL().RawQuery)
		if err != nil {
			log.Printf("[Socket] ❌ Error parsing URL query for socket %s: %v", s.ID(), err)
			return err
		}
		userId := urlValues.Get("userId")

		// // 2. Lấy token từ payload 'auth'
		// // Client sẽ gửi nó qua header 'X-Socketio-Auth' dưới dạng một chuỗi JSON
		// authHeader := s.RemoteHeader().Get("X-Socketio-Auth")

		// log.Printf("[Socket] 🔎 Handshake data: userId='%s', AuthHeader='%s'", userId, authHeader)

		// var authPayload struct {
		// 	Token string `json:"token"`
		// }

		// var token string
		// if authHeader != "" {
		// 	// Giải mã chuỗi JSON từ header
		// 	err := json.Unmarshal([]byte(authHeader), &authPayload)
		// 	if err != nil {
		// 		log.Printf("[Socket] 🚫 Connection rejected for socket %s: Invalid auth payload format.", s.ID())
		// 		return fmt.Errorf("authentication failed: invalid auth format")
		// 	}
		// 	token = authPayload.Token
		// }

		// 3. Kiểm tra dữ liệu
		// if userId == "" || token == "" {
		if userId == "" {
			log.Printf("[Socket] 🚫 Connection rejected for socket %s: Missing userId or token.", s.ID())
			return fmt.Errorf("authentication failed: missing credentials")
		}

		log.Printf("[Socket] 👍 Connection accepted for user %s with socket %s", userId, s.ID())

		if userId != "" {
			userSocketMap[userId] = s.ID()
			log.Println("User mapped: ", userId, "->", s.ID())
		}
		time.Sleep(1000 * time.Millisecond) // chờ 0.5s để frontend subscribe listener
		//Broadcast the updated list of online users to all clients.
		onlineUsers := make([]string, 0, len(userSocketMap))
		for k := range userSocketMap {
			onlineUsers = append(onlineUsers, k)
		}
		log.Println("[TESTING RIGHT HERE: ] Online users:", onlineUsers)
		// s.Emit("getOnlineUsers", onlineUsers)
		// server.BroadcastToRoom("/socket.io", s.ID(), "getOnlineUsers", onlineUsers)
		server.BroadcastToNamespace("/socket.io/", "getOnlineUsers", onlineUsers)
		log.Println("[BROADCAST: ] Successfully broadcasted:", onlineUsers)
		// Liệt kê tất cả namespace hiện tại
		log.Println("===== Active namespaces =====")
		log.Println("Socket connected to namespace:", s.Namespace())
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

		server.BroadcastToNamespace("/socket.io/", "getOnlineUsers", onlineUsers)
	})

	return server
}

// ErrServerClosed
type ErrServerClosed struct{}

func (e *ErrServerClosed) Error() string {
	return "server closed"
}

func EmitToSocket(userId string, event string, data interface{}) {
	if sid, ok := userSocketMap[userId]; ok {
		// server.BroadcastToRoom("/", sid, event, data)
		server.ForEach("/socket.io/", "", func(s socketio.Conn) {
			if s.ID() == sid {
				s.Emit(event, data)
			}
		})
	}
}

//
// GetSocketManager provides a single instance of the SocketManager using a singleton pattern.
// func GetSocketManager() *SocketManager {
// 	once.Do(func() {
// 		server := socketio.NewServer(&engineio.Options{
// 			Transports: []transport.Transport{
// 				polling.Default,
// 				websocket.Default,
// 			},
// 		})

// 		instanceSM = &SocketManager{
// 			server:        server,
// 			userSocketMap: make(map[string]string),
// 		}

// 		instanceSM.setupEventHandlers()
// 	})
// 	return instanceSM
// }

// // setupEventHandlers configures the socket connection logic.
// func (sm *SocketManager) setupEventHandlers() {
// 	sm.server.OnConnect("/", func(s socketio.Conn) error {
// 		log.Println("A user connected:", s.ID())

// 		// Use a middleware or an auth check here to validate the user.
// 		// For simplicity, we'll continue using the query parameter.
// 		urlValues, err := url.ParseQuery(s.URL().RawQuery)
// 		if err != nil {
// 			log.Printf("Error parsing URL query for socket %s: %v", s.ID(), err)
// 			s.Close()
// 			return err
// 		}
// 		userId := urlValues.Get("userId")

// 		if userId != "" {
// 			sm.addUser(userId, s.ID())
// 		}

// 		return nil
// 	})

// 	sm.server.OnDisconnect("/", func(s socketio.Conn, reason string) {
// 		log.Println("A user disconnected:", s.ID(), "Reason:", reason)
// 		sm.removeUser(s, s.ID())
// 	})

// 	sm.server.OnError("/", func(s socketio.Conn, e error) {
// 		log.Println("Socket error:", e)
// 	})
// }

// // addUser safely adds a user to the online list and broadcasts the update.
// func (sm *SocketManager) addUser(userId, socketId string) {
// 	sm.mu.Lock()
// 	defer sm.mu.Unlock()

// 	// Add or update the user's socket ID
// 	sm.userSocketMap[userId] = socketId
// 	log.Println("User mapped:", userId, "->", socketId)

// 	// Broadcast the updated list
// 	sm.broadcastOnlineUsersLocked()
// }

// // removeUser safely removes a user and broadcasts the update.
// func (sm *SocketManager) removeUser(s socketio.Conn, socketId string) {
// 	sm.mu.Lock()
// 	defer sm.mu.Unlock()

// 	var userIdToRemove string
// 	for id, sid := range sm.userSocketMap {
// 		if sid == socketId {
// 			userIdToRemove = id
// 			break
// 		}
// 	}

// 	if userIdToRemove != "" {
// 		delete(sm.userSocketMap, userIdToRemove)
// 		log.Println("User unmapped:", userIdToRemove)
// 	}

// 	// Broadcast the updated list
// 	sm.broadcastOnlineUsersLocked()
// }

// // broadcastOnlineUsersLocked is a helper that broadcasts the list.
// // It assumes the mutex is already locked.
// func (sm *SocketManager) broadcastOnlineUsersLocked() {
// 	onlineUsers := make([]string, 0, len(sm.userSocketMap))
// 	for k := range sm.userSocketMap {
// 		onlineUsers = append(onlineUsers, k)
// 	}

// 	// s.Emit("getOnlineUsers", onlineUsers)

// 	log.Println("Broadcasting online users:", onlineUsers)
// 	sm.server.BroadcastToNamespace("/socket.io", "getOnlineUsers", onlineUsers)
// }

// // EmitToSocket sends an event to a specific user.
// func (sm *SocketManager) EmitToSocket(userId string, event string, data interface{}) {
// 	sm.mu.RLock()
// 	sid, ok := sm.userSocketMap[userId]
// 	sm.mu.RUnlock()

// 	if ok {
// 		sm.server.BroadcastToRoom("/", sid, event, data)
// 	}
// }

// // Server returns the underlying socket.io server instance.
// func (sm *SocketManager) Server() *socketio.Server {
// 	return sm.server
// }

// // GetReceiverSocketId returns the socket ID for a given user ID.
// func (sm *SocketManager) GetReceiverSocketId(userId string) string {
// 	sm.mu.RLock()
// 	defer sm.mu.RUnlock()
// 	return sm.userSocketMap[userId]
// }
