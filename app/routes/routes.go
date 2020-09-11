package routes

import (
	"todo-api/app/controllers"

	"github.com/gin-gonic/gin"
)

// InitializeRoutes ...
func InitializeRoutes(router *gin.Engine) {

	v1 := router.Group("/api/v1")

	//Users routes
	v1.POST("/register", controllers.CreateUser)
	//v1.POST("/address", controllers.CreateAddress)
	//v1.POST("/login", controllers.Login)
	//v1.GET("/profile/:id", controllers.GetUserByID)

	// Address
	// v1.POST("/address", app.TokenAuthMiddleware(), controllers.CreateAddress)

	// Task
	//v1.POST("/task", middlewares.TokenAuthMiddleware(), controllers.AddTask)
	//v1.GET("/task", middlewares.TokenAuthMiddleware(), controllers.GetAllTask)
	//v1.PUT("/task/:id", middlewares.TokenAuthMiddleware(), controllers.UpdateTask)
	//v1.DELETE("/task/:id", middlewares.TokenAuthMiddleware(), controllers.DeleteTask)
}
