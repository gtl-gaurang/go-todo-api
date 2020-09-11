package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
)

// DataSource ..
type DataSource struct {
	DB *gorm.DB
}

var DB = &DataSource{}

// Initialize ... Used for initialize the DB connection
func Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) *DataSource {
	var err error

	fmt.Println("Database")
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser,
			DbPassword, DbHost, DbPort, DbName)
		DB.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		DB.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error connecting to postgres:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else {
		fmt.Println("Unknown Driver")
	}

	DB.DB.Debug().AutoMigrate(
		&User{},
		&UserAddress{},
		&Task{},
	)

	// Add Foreign Key
	DB.DB.Model(&UserAddress{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	DB.DB.Model(&Task{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	return DB
}
