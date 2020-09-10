package app

// InitializeRoutes ...
func (app *App) InitializeRoutes() {

	v1 := app.router.Group("/api/v1")

	//Users routes
	v1.POST("/register", app.CreateUser())
	//v1.POST("/address", controllers.CreateAddress)
	//v1.POST("/login", app.Login)
	//v1.GET("/profile/:id", app.GetUserByID)

	// Address
	//v1.POST("/address", middlewares.TokenAuthMiddleware(), app.CreateAddress)

	// Task
	//v1.POST("/task", middlewares.TokenAuthMiddleware(), app.AddTask)
	//v1.GET("/task", middlewares.TokenAuthMiddleware(), app.GetAllTask)
	//v1.PUT("/task/:id", middlewares.TokenAuthMiddleware(), app.UpdateTask)
	//v1.DELETE("/task/:id", middlewares.TokenAuthMiddleware(), app.DeleteTask)
}
