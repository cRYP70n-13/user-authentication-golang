package routes

import (
	"log"
	gin "github.com/gin-gonic/gin"
)

router := gin.Default()
.POST("/questions", CreateQuestionEndpoint)

if err r.Run(":8080"); err != nil {
	log.Fatal(err)
}
