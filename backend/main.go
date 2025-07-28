package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	authToken "github.com/tph-kds/chat_realtime/backend/internal/api/handlers"
	"github.com/tph-kds/chat_realtime/backend/internal/api/routes"
	"github.com/tph-kds/chat_realtime/backend/internal/configs"
	"github.com/tph-kds/chat_realtime/backend/internal/database"
)

func main() {

	// Connect to MongoDB
	log.Println("Connecting to MongoDB...")
	// Load environment variables
	err := configs.LoadConfigEnv("backend/.env")
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Read MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

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

	routes.SetupRoutes(r)

	// Set up the Port Server
	configs.SetServerPort(configs.PORT)

	//Start the server
	r.Run(":" + configs.PORT)
	log.Println("Server is running on port " + configs.PORT)
}
