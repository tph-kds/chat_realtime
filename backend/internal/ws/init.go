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
	return socketio.NewServer(nil)
}

func ConnectSocketServer(s *socketio.Server) {
	socketServer = s
}

func DisconnectSocketServer() {
	socketServer = nil
}
