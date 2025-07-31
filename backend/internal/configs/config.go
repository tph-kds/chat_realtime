package configs

import (
	"log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func LoadConfigEnv(file_name string) error {
	err := godotenv.Load(file_name)
	if err != nil {
		log.Println("Error loading .env file or file not found:", err)
		return err
	}
	log.Println("Environment variables loaded successfully from", file_name)
	return nil
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.StandardClaims
}

type ServerConfig struct {
	Port string `json:"port"`
}

var serverConfig ServerConfig

func GetServerPort() string {
	return serverConfig.Port
}

func SetServerPort(port string) {
	serverConfig.Port = port
}

// UpdateProfile Function
type UpdateProfileRequest struct {
	ProfilePic string `json:"profile_pic" binding:"required"`
}

type User struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone" binding:"required"`
	ProfilePic string `json:"profile_pic" binding:"required"`
}
