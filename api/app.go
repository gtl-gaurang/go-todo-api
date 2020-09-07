package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
)

// App ... Application
type App struct {
	Router *gin.Engine
}

var errList = make(map[string]string)
var app = App{}

//Run ...
func Run() {

	app.Router = gin.Default()
	app.Router.Use(CORSMiddleware())
	app.InitializeRoutes()

	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", apiPort)
	app.Run(apiPort)
}

//Run ... Handle the Http request
func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
