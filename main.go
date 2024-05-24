package main

import (
	"ArautosApp/controllers"
	"ArautosApp/initializers"
	"ArautosApp/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/cadastro", controllers.SingUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.ValidateLogin)

	r.Run()
}
