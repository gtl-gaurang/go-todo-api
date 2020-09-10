package app

import (
	"todo-api/app/controllers"

	"github.com/gin-gonic/gin"
)

//CreateUser ...
func (app *App) CreateUser() func(*gin.Context) {
	return controllers.CreateUser
}
