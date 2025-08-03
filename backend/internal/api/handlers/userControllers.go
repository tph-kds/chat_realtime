package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/tph-kds/chat_realtime/backend/internal/configs"
	"github.com/tph-kds/chat_realtime/backend/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var user models.User

		// Get User Input
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		// Validate User Input
		if validateErr := validate.Struct(user); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error(), "message": "Invalid user input"})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{
			"$or": []bson.M{
				{"email": user.Email},
				{"phone": user.Phone},
			},
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to check existing user"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists with this email or phone number"})
			return
		}
		// Generate rest of the user data
		user.Password = HashPassword(user.Password)
		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		accessToken, refreshToken := GenerateTokens(*user.Email, user.User_id, *user.Role)
		user.Token = &accessToken
		user.Refresh_token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error(), "message": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		//Get USer Input
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
			return
		}
		// Verify password
		passwordIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}
		token, refreshToken := GenerateTokens(*foundUser.Email, *&foundUser.User_id, *foundUser.Role)
		UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, gin.H{
			"user":          foundUser,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func LogOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedUserId := c.Param("id")
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

		tokenUserId := tokenClaims.UserID
		userType := tokenClaims.Role

		if (userType != "ADMIN" && userType != "admin") && tokenUserId != requestedUserId {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this user"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": requestedUserId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve claims from context
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

		// Find all users
		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var users []models.User
		if err := cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the list of users
		c.JSON(http.StatusOK, users)
	}
}

func UpdateProfileUser() gin.HandlerFunc {
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
		var reqBody configs.UpdateProfileRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid request body in update profile"})
			return
		}

		// Validate User Input
		if reqBody.ProfilePic == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Profile pic is required", "message": "Invalid request body in update profile"})
			return
		}

		// Upload profile pic to cloudinary
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		uploadResutl, err := GetCloudinary().Upload.Upload(ctx, reqBody.ProfilePic, uploader.UploadParams{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to upload profile pic"})
			return
		}

		// Update into database
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var updatedUser configs.User
		err = userCollection.FindOneAndUpdate(
			ctx,
			bson.M{"user_id": userId},
			bson.M{"$set": bson.M{"profile_pic": uploadResutl.SecureURL}},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&updatedUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to update profile pic in database"})
			return
		}

		// Send successful response
		c.JSON(http.StatusOK, updatedUser)

	}
}

func CheckUser() gin.HandlerFunc {
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

		// Get user from database
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var user configs.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to get user from database"})
			return
		}

		// 2. Send success response with user data
		c.JSON(http.StatusOK, user)

	}
}
