package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Make Type -> enum
type Category struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name" binding:"required"`
	Type  string             `json:"type" binding:"required"`
	Owner primitive.ObjectID `bson:"owner,omitempty"`
}
