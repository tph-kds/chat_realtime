package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/tph-kds/chat_realtime/backend/internal/api/handlers"
	"github.com/tph-kds/chat_realtime/backend/internal/api/middleware"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/signup", controllers.SignUp())
	router.POST("/login", controllers.Login())

	protected := router.Group("/")

	protected.Use(middleware.Authenticate())
	{
		protected.GET("/users", controllers.GetUsers())
		protected.GET("/users/:id", controllers.GetUser())
	}
}
