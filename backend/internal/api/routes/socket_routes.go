package routes

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func SetUpSocketRoutes(r *gin.Engine, server *socketio.Server) {

	// Mount Socket.IO into Gin
	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))
}
