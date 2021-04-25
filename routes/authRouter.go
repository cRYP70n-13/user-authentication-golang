package routes

import (
	controller "user-athentication-golang/controllers"

	"github.com/gin-gonic/gin"
)

//UserRoutes function
func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/v1/users/signup", controller.SignUp())
	incomingRoutes.POST("/api/v1/users/login", controller.Login())
}
