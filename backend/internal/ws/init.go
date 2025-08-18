package ws

import (
	socketio "github.com/googollee/go-socket.io"
)

// ======================== Socket Initialization ==========================

var socketServer *socketio.Server = nil

func SetSocketServer(s *socketio.Server) {
	socketServer = s
}

func GetSocketServer() *socketio.Server {
	return socketServer
}

func NewSocketServer() *socketio.Server {
	return NewServerCustom()
}

func ConnectSocketServer(s *socketio.Server) {
	socketServer = s
}

func DisconnectSocketServer() {
	socketServer = nil
}

// type SocketManager struct {
// 	server        *socketio.Server
// 	userSocketMap map[string]string
// 	mu            sync.RWMutex // RWMUTEX for safe concurency access to the map
// }

// var (
// 	instanceSM *SocketManager
// 	once       sync.Once //Use it to ensure that the instance is created only once
// )
