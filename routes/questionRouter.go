package routes

import (
	"context"
	"log"
	"user-athentication-golang/controllers"
	"user-athentication-golang/database"
	"user-athentication-golang/middleware"

	gin "github.com/gin-gonic/gin"
)

func QuestionRoutes(incomingRoutes *gin.Engine) {
	ctx := context.Background()
	elasticClient, err := database.GetESClient()
	if err != nil {
		log.Println("There is an error somewhere", err.Error())
	}
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/questions", controllers.GetAll(ctx, elasticClient))
}
