package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func SetUpSocketRoutes(r *gin.Engine, server *socketio.Server) {
	// // Wrap with CORS allowing your front-end (e.g., localhost:5173)
	// corsHandler := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:5173"},
	// 	AllowCredentials: true,
	// 	AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// })
	// // Mount Socket.IO into Gin
	// r.GET("/socket.io/*any", gin.WrapH(corsHandler.Handler(server)))
	// r.POST("/socket.io/*any", gin.WrapH(corsHandler.Handler(server)))
	log.Println("🔌 Mounting /socket.io/*any routes...")
	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))
}
