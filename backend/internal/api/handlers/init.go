package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/tph-kds/chat_realtime/backend/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtKey []byte

func SetJWTKey(key string) {
	jwtKey = []byte(key)
}

func GetJWTKey() []byte {
	return []byte(jwtKey)
}

var client *mongo.Client

func SetClient(c *mongo.Client) {
	client = c
}

var dbName string

func SetDBName(name string) {
	dbName = name
}

var collectionName string

func SetCollectionName(name string) {
	collectionName = name
}

// var validate = validator.New()
var validate *validator.Validate

func InitValidator() {
	if validate == nil {
		validate = validator.New()
	}
	// validate = validate
}

var userCollection *mongo.Collection

func SetUserCollection(c *mongo.Client) error {
	if client == nil {
		client = c
	}
	if dbName == "" {
		dbName = "chat_realtime"
	}
	// Initialize the user collection
	userCollection = database.OpenCollection(client, dbName, collectionName)

	return nil
}
