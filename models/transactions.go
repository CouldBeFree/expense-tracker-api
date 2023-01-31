package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Make Type -> enum
type Transaction struct {
	ID       primitive.ObjectID       `json:"id" bson:"_id"`
	Category primitive.ObjectID       `bson:"category,omitempty" json:"category"`
	Amount   string                   `json:"amount" binding:"required"`
	Owner    primitive.ObjectID       `bson:"owner,omitempty" json:"owner"`
	InvDt    primitive.DateTime       `bson:"invdt,omitempty" json:"invdt,omitempty"`
	Date     string                   `json:"date" binding:"required"`
	Cat      []map[string]interface{} `json:"cat" bson:"cat"`
}
