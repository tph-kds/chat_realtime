package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	First_name    *string            `bson:"first_name" json:"first_name" validate:"required,min=2,max=20"`
	Last_name     *string            `bson:"last_name" json:"last_name" validate:"required,min=2,max=20"`
	Password      *string            `bson:"password" json:"password" validate:"required,min=6"`
	Email         *string            `bson:"email" json:"email" validate:"email,required"`
	Phone         *string            `bson:"phone" json:"phone" validate:"required"`
	Token         *string            `bson:"token" json:"token,omitempty"`
	Role          *string            `bson:"role" json:"role" validate:"required,eq=admin|eq=user"`
	Refresh_token *string            `bson:"refresh_token" json:"refresh_token,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	User_id       string             `bson:"user_id" json:"user_id"`
	ProfilePic    string             `bson:"profile_pic" json:"profile_pic"`
}

type UserTest struct {
	Email      string `bson:"email" json:"email" validate:"email,required"`
	FullName   string `bson:"fullName" json:"fullName" validate:"required"`
	Password   string `bson:"password" json:"password" validate:"required,min=6"`
	ProfilePic string `bson:"profile_pic" json:"profile_pic" validate:"required"`
}

type Message struct {
	// ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SenderId   primitive.ObjectID `bson:"senderId" json:"senderId"`
	ReceiverId primitive.ObjectID `bson:"receiverId" json:"receiverId"`
	Text       string             `bson:"text" json:"text"`
	Image      string             `bson:"image,omitempty" json:"image,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
