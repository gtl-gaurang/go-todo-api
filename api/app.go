package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo-api/api/models"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
	"github.com/joho/godotenv"
)

// App ... Application
type App struct {
	router *gin.Engine
	db     *models.DataSource
}

var errList = make(map[string]string)

//Run ...
func Run() {
	var app = App{}
	app.router = gin.Default()
	app.router.Use(CORSMiddleware())
	app.InitializeRoutes()

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting values")
	}

	app.db.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", apiPort)
	app.Run(apiPort)
}

//Run ... Handle the Http request
func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.router))
}
