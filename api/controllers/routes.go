package controllers

import "todo-api/api/middlewares"

func (s *Server) initializeRoutes() {

	v1 := s.Router.Group("/api/v1")

	//Users routes
	v1.POST("/register", s.CreateUser)
	v1.POST("/login", s.Login)
	v1.GET("/profile/:id", s.GetUserByID)

	v1.POST("/task", middlewares.TokenAuthMiddleware(), s.AddTask)
	v1.GET("/task", middlewares.TokenAuthMiddleware(), s.GetAllTask)
	v1.PUT("/task/:id", middlewares.TokenAuthMiddleware(), s.UpdateTask)
	v1.DELETE("/task/:id", middlewares.TokenAuthMiddleware(), s.DeleteTask)
}
