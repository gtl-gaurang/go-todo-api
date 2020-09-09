package api

import (
	"todo-api/api/controllers"

	"github.com/gin-gonic/gin"
)

//CreateUser ...
func (a *App) CreateUser() func(*gin.Context) {
	return controllers.CreateUser
}
