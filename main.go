package main

import (
	_ "fmt"

	handlers "expense-tracker-api/handlers"

	"github.com/gin-gonic/gin"
)

var authHandler *handlers.AuthHandler

func main() {
	router := gin.Default()
	router.POST("/register", authHandler.RegisterUser)
	router.Run(":5050")
}
