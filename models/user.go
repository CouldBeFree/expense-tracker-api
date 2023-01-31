package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID   `bson:"_id"`
	Username     string               `json:"username" binding:"required"`
	Password     string               `json:"password" binding:"required"`
	Email        string               `json:"email" binding:"required,email"`
	Categories   []primitive.ObjectID `bson:"categories,omitempty"`
	Transactions []primitive.ObjectID `bson:"transactions,omitempty"`
}

type LogggedInUser struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
