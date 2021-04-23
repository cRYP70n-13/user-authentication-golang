package main

import (
	"context"
	"fmt"
	"os"

	"user-athentication-golang/controllers"
	"user-athentication-golang/database"
	routes "user-athentication-golang/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

const (
	indexName = "question_index"
)

func main() {
	port := os.Getenv("PORT")
	ctx := context.Background()
	elasticClient, err := database.GetESClient()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(-1)
	}

	// Create the index in elasticsearch
	controllers.CreateIndexIfDoesNotExists(ctx, elasticClient, indexName)

	// Insert the questions
	fmt.Println("------------Inserting Questions-------------")
	controllers.InsertQuestion(ctx, elasticClient)

	fmt.Println("------------GetAll-------------")
	controllers.GetAll(ctx, elasticClient)

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// API-1
	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	// API-2
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
