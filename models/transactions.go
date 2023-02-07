package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Make Type -> enum
// Hobbies: []string{"IT","Travel"}
type Transaction struct {
	ID           primitive.ObjectID       `json:"id" bson:"_id"`
	Category     primitive.ObjectID       `bson:"category,omitempty" json:"category"`
	Amount       int                      `json:"amount" binding:"required"`
	Converted    int                      `bson:"converted,omitempty" json:"converted,omitempty"`
	Owner        primitive.ObjectID       `bson:"owner,omitempty" json:"owner"`
	InvDt        primitive.DateTime       `bson:"invdt,omitempty" json:"invdt,omitempty"`
	Date         string                   `json:"date" binding:"required"`
	Cat          []map[string]interface{} `json:"cat" bson:"cat"`
	Transactions []map[string]interface{} `json:"transactions" bson:"transactions"`
}

// Make Type -> enum
// Hobbies: []string{"IT","Travel"}
type TransactionCategory struct {
	Total int                      `json:"total" bson:"total"`
	Cat   []map[string]interface{} `json:"category" bson:"category"`
	ID    primitive.ObjectID       `json:"id" bson:"_id"`
}
