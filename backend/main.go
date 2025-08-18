package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/gofiber/fiber/v2/middleware/cors"
	// "time"

	authToken "github.com/tph-kds/chat_realtime/backend/internal/api/handlers"
	"github.com/tph-kds/chat_realtime/backend/internal/api/routes"
	"github.com/tph-kds/chat_realtime/backend/internal/configs"
	"github.com/tph-kds/chat_realtime/backend/internal/database"
	"github.com/tph-kds/chat_realtime/backend/internal/ws"
)

func main() {

	// Connect to MongoDB
	log.Println("Connecting to MongoDB...")
	// Load environment variables
	err := configs.LoadConfigEnv("backend/.env")
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// // Get The initialize Args
	// Read MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}
	// Get Cloudinary URL from environment
	cldURL := os.Getenv("CLOUDINARY_URL")
	if cldURL == "" {
		log.Fatal("CLOUDINARY_URL environment variable is not set")
	}
	// Initialize Cloudinary
	cldClient, err := database.ConnectCloudinary(cldURL)
	if err != nil {
		log.Fatalf("Cloudinary initialization error: %v", err)
	}

	authToken.SetCloudinary(cldClient)

	// Connect to MongoDB
	client, err := database.ConnectDB(mongoURI)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	authToken.SetClient(client)
	// Set the database name
	authToken.SetDBName(configs.DB_NAME)
	// Set the collection name
	authToken.SetCollectionName(configs.COLLECTION_NAME)
	// Set Validator
	authToken.InitValidator()

	// Set User Collection
	authToken.SetUserCollection(client)

	// Set JWT Key
	jwtKey := configs.GenerateRandomKey()

	authToken.SetJWTKey(jwtKey)

	// Initialize MongoDB connection and gin router
	r := gin.Default()

	//  Setup CORS

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	}))

	// Set up Websocket
	socketServer := ws.NewSocketServer()
	ws.SetSocketServer(socketServer)

	// // Get the socket manager instance
	// socketManager := ws.GetSocketManager()

	// // Start the socket server in a goroutine
	// go func() {
	// 	if err := socketManager.Server().Serve(); err != nil {
	// 		log.Fatalf("socketio listen error: %s\n", err)
	// 	}
	// }()
	// defer socketManager.Server().Close()
	// // Set up Socket.IO URL

	go func() {
		if err := socketServer.Serve(); err != nil {
			log.Fatalf("Error starting WebSocket server: %v", err)
		}
	}()
	log.Println("✅ Socket.IO server started")
	defer socketServer.Close()

	routes.SetUpSocketRoutes(r, socketServer)
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/socket.io") {
			log.Printf("[GIN-DEBUG] %s %s", c.Request.Method, c.Request.URL.String())
		}
		c.Next()
	})
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/socket.io") {
			log.Printf("[GIN-DEBUG] %s %s Upgrade: %s",
				c.Request.Method,
				c.Request.URL.String(),
				c.Request.Header.Get("Upgrade"),
			)
		}
		c.Next()
	})
	routes.SetupRoutes(r)

	// Set up the Port Server
	configs.SetServerPort(configs.PORT)

	//Start the server
	r.Run(":" + configs.PORT)
	log.Println("Server is running on port " + configs.PORT)
}
