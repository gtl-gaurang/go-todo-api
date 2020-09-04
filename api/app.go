package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo-api/api/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
	"github.com/joho/godotenv"
)

// App ... Application
type App struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)
var app = App{}

//Run ...
func Run() {
	var err error
	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting values")
	}

	app.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	// This is for testing, when done, do well to comment
	// seed.Load(server.DB)

	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", apiPort)

	app.Run(apiPort)

}

// Initialize ... Used for initialize the DB connection
func (app *App) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		app.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		app.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error connecting to postgres:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else {
		fmt.Println("Unknown Driver")
	}

	app.DB.Debug().AutoMigrate(
		&models.User{},
		&models.UserAddress{},
		&models.Task{},
	)

	// Add Foreign Key
	app.DB.Model(&models.UserAddress{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	app.DB.Model(&models.Task{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")

	app.Router = gin.Default()
	app.Router.Use(CORSMiddleware())
	app.InitializeRoutes()
}

//Run ... Handle the Http request
func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
