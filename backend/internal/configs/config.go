package configs

import (
	"log"
	"time"

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
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`

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
	ProfilePic string    `bson:"profile_pic" json:"profile_pic" binding:"required"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	// CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

type User struct {
	FirstName  string    `bson:"first_name" json:"first_name" binding:"required"`
	LastName   string    `bson:"last_name" json:"last_name" binding:"required"`
	Email      string    `bson:"email" json:"email" binding:"required,email"`
	Phone      string    `bson:"phone" json:"phone" binding:"required"`
	ProfilePic string    `bson:"profile_pic" json:"profile_pic" binding:"required"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}
