package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nimaibhat/GoCrudBookManagement/http"
	"github.com/nimaibhat/GoCrudBookManagement/repository"
	"github.com/nimaibhat/GoCrudBookManagement/middleware"
)

func main() {
	// Create a database connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Connect(context.TODO()); err != nil {
		log.Fatal(err)
	}

	// Create a repository
	repo := repository.NewRepository(client.Database("books"))

	// Create an HTTP server
	server := http.NewServer(repo)

	// Create a Gin router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Apply custom recovery middleware
	router.Use(middleware.RecoveryMiddleware())

	// Define routes
	router.POST("/books", server.CreateBook)
	router.GET("/books/:id", server.GetBook)
	router.PUT("/books/:id", server.UpdateBook)
	router.DELETE("/books/:id", server.DeleteBook)

	// Start the router
	router.Run(":8080")
}
