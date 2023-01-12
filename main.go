package main

import (
	"context"
	"log"

	handlers "expense-tracker-api/handlers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

var authHandler *handlers.AuthHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	collectionUsers := client.Database("expense").Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func main() {
	router := gin.Default()
	router.POST("/register", authHandler.RegisterUser)
	router.POST("/signin", authHandler.SignInHandler)
	router.Run(":5050")
}
