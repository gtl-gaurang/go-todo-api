package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo-api/app/models"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
	"github.com/joho/godotenv"
)

//Run ... Handle the Http request
func Run() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	InitializeRoutes(router)

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting values")
	}

	models.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", apiPort)
	log.Fatal(http.ListenAndServe(apiPort, router))
}
