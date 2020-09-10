package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
)

type DataSource struct {
	DB *gorm.DB
}

// Initialize ... Used for initialize the DB connection
func (db *DataSource) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	fmt.Println("Database")
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser,
			DbPassword, DbHost, DbPort, DbName)
		db.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		db.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error connecting to postgres:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else {
		fmt.Println("Unknown Driver")
	}

	db.DB.Debug().AutoMigrate(
		&User{},
		//&UserAddress{},
		//&Task{},
	)

	// Add Foreign Key
	//db.DB.Model(&UserAddress{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	//db.DB.Model(&Task{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
}
