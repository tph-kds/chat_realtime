package handlers

import (
	"context"
	// "fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tph-kds/chat_realtime/backend/internal/configs"
	"github.com/tph-kds/chat_realtime/backend/internal/domain/models"
	lib "github.com/tph-kds/chat_realtime/backend/internal/lib"
	wsCustom "github.com/tph-kds/chat_realtime/backend/internal/ws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get claims from context
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Type assertion to get the claims object
		tokenClaims, ok := claims.(*configs.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		// Check if the user is an ADMIN
		if tokenClaims.Role != "ADMIN" && tokenClaims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		userId := tokenClaims.UserID
		myId, _ := primitive.ObjectIDFromHex(userId)

		userToChatIdStr := c.Param("id")
		userToChatId, _ := primitive.ObjectIDFromHex(userToChatIdStr)

		// Get user from database
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// var user configs.User
		filter := bson.M{
			"$or": []bson.M{
				{"senderId": myId, "receiverId": userToChatId},
				{"senderId": userToChatId, "receiverId": myId},
			},
		}

		cursor, err := userCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var messages []models.Message
		if err := cursor.All(ctx, &messages); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to get messages"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"messages": messages})
	}
}

func GetUsersForSidebar() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get claims from context
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Type assertion to get the claims object
		tokenClaims, ok := claims.(*configs.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		// Check if the user is an ADMIN
		if tokenClaims.Role != "ADMIN" && tokenClaims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Proceed to get the users list from the database
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		loggedUserId, _ := primitive.ObjectIDFromHex(tokenClaims.UserID)
		filter := bson.M{"_id": bson.M{"$ne": loggedUserId}}

		cursor, err := userCollection.Find(ctx, filter, options.Find().SetProjection(bson.M{"password": 0}))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to get users"})
			return
		}
		var users []configs.User

		if err := cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to get users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})

	}
}

func SendMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get claims from context
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Type assertion to get the claims object
		tokenClaims, ok := claims.(*configs.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		// Check if the user is an ADMIN
		if tokenClaims.Role != "ADMIN" && tokenClaims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		senderId, _ := primitive.ObjectIDFromHex(tokenClaims.UserID)
		receiverId, _ := primitive.ObjectIDFromHex(c.Param("id"))

		// fmt.Println("Sender ID:", senderId, "Receiver ID:", receiverId)

		var body struct {
			Text  string `json:"text"`
			Image string `json:"image"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var imageUrl string
		if body.Image != "" {
			uploadResutl, err := lib.UploadToCloudinary(body.Image)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to upload image"})
				return
			}
			imageUrl = uploadResutl.SecureURL
		}

		newMessages := models.Message{
			SenderId:   senderId,
			ReceiverId: receiverId,
			Text:       body.Text,
			Image:      imageUrl,
			CreatedAt:  time.Now(),
		}

		_, err := userCollection.InsertOne(context.Background(), newMessages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to send message"})
			return
		}
		wsCustom.EmitToSocket("newMessage", receiverId.Hex(), newMessages)
		c.JSON(http.StatusOK, gin.H{"status": "Message sent successfully", "message": newMessages})

	}
}
