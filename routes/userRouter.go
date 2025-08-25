package routes

import (
	"github.com/0xk4n3ki/OAuth2.0-golang/controllers"
	"github.com/0xk4n3ki/OAuth2.0-golang/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
}
