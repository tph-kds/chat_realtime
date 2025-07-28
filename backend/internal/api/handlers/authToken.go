package handlers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tph-kds/chat_realtime/backend/internal/configs"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func ValidateToken(tokenString string) (*configs.Claims, error) {
	// Using dynamically set JWT Key here
	secretKey := GetJWTKey()
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &configs.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	// Check if the token is valid and return the claims
	if claims, ok := token.Claims.(*configs.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Invalid token validation function. Please implement it.")
}

func GenerateTokens(email, userID, userType string) (string, string) {
	log.Printf("JWT Key %v Type %T", jwtKey, jwtKey)

	// Token expirations times
	tokenExpiry := time.Now().Add(24 * time.Hour).Unix()
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour).Unix()

	claims := &configs.Claims{
		UserID: userID,
		Email:  email,
		Role:   userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiry,
		},
	}

	refreshClaims := &configs.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiry,
		},
	}

	// Generate tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessToken, err := accessToken.SignedString(jwtKey)
	if err != nil {
		log.Fatalf("Error signing access token: %v", err)
		panic("Failed to sign access token")
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		log.Fatalf("Error signing refresh token: %v", err)
		panic("Failed to sign refresh token")
	}

	return signedAccessToken, signedRefreshToken
}

func HashPassword(password *string) *string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	hashedPassword := string(bytes)
	return &hashedPassword
}

func UpdateAllTokens(signedToken, refreshToken, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// userCollection := database.OpenCollection(client, dbName, "users")

	// Create a update object
	updateObj := bson.D{
		{"$set", bson.D{
			{"token", signedToken},
			{"refresh_token", refreshToken},
			{"updated_at", time.Now()},
		},
		},
	}

	//Create a filter
	filter := bson.M{"user_id": userID}
	_, err := userCollection.UpdateOne(ctx, filter, updateObj)

	return err

}

func VerifyPassword(hashedPassword, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, "Invalid email or password"
}
