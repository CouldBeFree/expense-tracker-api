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
var categoriesHandler *handlers.CategoryHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	collectionUsers := client.Database("expense").Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)

	collectionCategories := client.Database("expense").Collection("categories")
	categoriesHandler = handlers.NewCategoryHandler(ctx, collectionCategories)
}

func main() {
	router := gin.Default()

	router.Use(authHandler.CORSMiddleware())

	router.POST("/register", authHandler.RegisterUser)
	router.POST("/signin", authHandler.SignInHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())

	{
		authorized.GET("/categories", categoriesHandler.ListCategory)
		authorized.POST("/create-category", categoriesHandler.CreateCategory)
		authorized.GET("/category/:id", categoriesHandler.GetCategory)
		authorized.DELETE("/category/:id", categoriesHandler.DeleteCategory)
		authorized.PUT("/category/:id", categoriesHandler.UpdateCategory)
	}

	router.Run(":5050")
}
