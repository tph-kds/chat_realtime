package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	First_name    *string            `json:"first_name" validate:"required,min=2,max=20"`
	Last_name     *string            `json:"last_name" validate:"required,min=2,max=20"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token,omitempty"`
	Role          *string            `json:"role" validate:"required,eq=admin|eq=user"`
	Refresh_token *string            `json:"refresh_token,omitempty"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}

type UserTest struct {
	Email      string `json:"email" validate:"email,required"`
	FullName   string `json:"fullName" validate:"required"`
	Password   string `json:"password" validate:"required,min=6"`
	ProfilePic string `json:"profilePic" validate:"required"`
}
