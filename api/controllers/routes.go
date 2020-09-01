package controllers

func (s *Server) initializeRoutes() {

	v1 := s.Router.Group("/api/v1")

	v1.POST("/task", s.AddTask)
	v1.GET("/task", s.GetAllTask)
	v1.PUT("/task/:id", s.UpdateTask)
	v1.DELETE("/task/:id", s.DeleteTask)
}
