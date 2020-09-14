package app

import (
	"todo-api/app/controllers"

	"github.com/gin-gonic/gin"
)

// InitializeRoutes ...
func InitializeRoutes(router *gin.Engine) {

	v1 := router.Group("/api/v1")

	//Users routes
	v1.POST("/register", controllers.CreateUser)
	v1.POST("/login", controllers.Login)
	v1.GET("/profile/:id", controllers.GetUserByID)

	// Address
	v1.POST("/address", TokenAuthMiddleware(), controllers.CreateAddress)

	// Task
	v1.POST("/task", TokenAuthMiddleware(), controllers.AddTask)
	v1.GET("/task", TokenAuthMiddleware(), controllers.GetAllTask)
	v1.PUT("/task/:id", TokenAuthMiddleware(), controllers.UpdateTask)
	v1.DELETE("/task/:id", TokenAuthMiddleware(), controllers.DeleteTask)
}
