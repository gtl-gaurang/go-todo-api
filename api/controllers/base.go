package controllers

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/api/middlewares"
	"todo-api/api/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
)

// Server ... created for handling the DB and Router
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)

// Initialize ... Used for initialize the DB connection
func (s *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		s.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		s.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error connecting to postgres:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else {
		fmt.Println("Unknown Driver")
	}

	s.DB.Debug().AutoMigrate(
		&models.Task{},
		&models.User{},
	)

	s.Router = gin.Default()
	s.Router.Use(middlewares.CORSMiddleware())
	s.initializeRoutes()
}

//Run ... Handle the Http request
func (s *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
