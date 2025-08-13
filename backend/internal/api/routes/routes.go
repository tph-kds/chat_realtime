package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/tph-kds/chat_realtime/backend/internal/api/handlers"
	"github.com/tph-kds/chat_realtime/backend/internal/api/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Check Auth
	router.POST("/signup", controllers.SignUp())
	router.POST("/logout", controllers.LogOut())
	router.POST("/login", controllers.Login())

	protected := router.Group("/")

	protected.Use(middleware.Authenticate())
	{
		protected.GET("/auth/check-auth", controllers.CheckAuth)
		protected.GET("/users", controllers.GetUsers())
		protected.GET("/users/:id", controllers.GetUser())
		protected.PUT("/users/:id/update-profile", controllers.UpdateProfileUser())
		protected.GET("/users/:id/check-user", controllers.CheckUser())

		// Websocket routes
		protected.GET("/online-users", controllers.GetUsersForSidebar())
		protected.GET("/messages/:id", controllers.GetMessages())
		protected.POST("/messages/send/:id", controllers.SendMessage())
	}
}
