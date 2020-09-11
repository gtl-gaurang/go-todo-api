package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
)

// DataRepo ...
type DataRepo interface {
	AddTask(t *Task) (*Task, error)
	GetAllTask(t *Task) (*[]Task, error)
	UpdateTask(tid uint32, t *Task) (*Task, error)
	FindTaskByID(tid uint32, t *Task) (*Task, error)
	DeleteTask(tid uint32, t *Task) (int64, error)

	AddUser(u *User) (*User, error)
	FindUserByID(uid uint32, u *User) (*User, error)
	GetAllUser() (*[]User, error)
	UpdateUser(u *User, uid uint32) (*User, error)
	DeleteUser(uid uint32) (int64, error)
	SignIn(email, password string) (map[string]interface{}, error)

	AddAddress(ua *UserAddress) (*UserAddress, error)
	GetAllAddress(ua *UserAddress) (*[]UserAddress, error)
	UpdateAddress(aid uint32, ua *UserAddress) (*UserAddress, error)
	DeleteAddress(aid uint32, ua *UserAddress) (int64, error)
}

// DataSource ..
type DataSource struct {
	*gorm.DB
}

// DB ..
var DB DataRepo

// Initialize ... Used for initialize the DB connection
func Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	var db *gorm.DB

	fmt.Println("Database")
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser,
			DbPassword, DbHost, DbPort, DbName)
		db, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		db, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error connecting to postgres:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else {
		fmt.Println("Unknown Driver")
	}

	db.Debug().AutoMigrate(
		&User{},
		&UserAddress{},
		&Task{},
	)

	// Add Foreign Key
	db.Model(&UserAddress{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	db.Model(&Task{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	DB = &DataSource{db}
}
