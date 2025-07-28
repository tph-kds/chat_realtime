package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(MONGODB_URI string) (*mongo.Client, error) {
	log.Println("Connecting to MongoDB...")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	log.Println("Connected to MongoDB successfully!")
	return client, nil
}

// var Client *mongo.Client = ConnectDB()

func OpenCollection(client *mongo.Client, dbName string, collectionName string) *mongo.Collection {
	if client == nil {
		log.Fatal("MongoDB client is not initialized. Please check the connection.")
	}
	return client.Database(dbName).Collection(collectionName)
}
