package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"
	"todo-api/api/auth"
	"todo-api/api/security"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//User ...
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	DOB       time.Time `gorm:"size:10;null" json:"dob"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//BeforeSave ...
func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

//Prepare ... Take data double check
func (u *User) Prepare() {
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Validate ...
func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "register":
		if u.Name == "" {
			err = errors.New("Required User full Name")
			errorMessages["required_name"] = err.Error()
		}
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	case "login":
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	default:
		if u.Name == "" {
			err = errors.New("Required User full Name")
			errorMessages["required_name"] = err.Error()
		}
	}

	return errorMessages
}

//AddUser ...
func (db *DataSource) AddUser(u *User) (*User, error) {

	fmt.Println("User value => ", &u)
	var err error
	err = db.DB.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

//FindUserByID ...
func (u *User) FindUserByID(uid uint32) (*User, error) {
	var err error
	err = DB.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

//GetAllUser ... Get All user
func (u *User) GetAllUser() (*[]User, error) {
	var err error
	tasks := []User{}
	err = DB.Debug().Model(&User{}).Limit(100).Find(&tasks).Error
	if err != nil {
		return &[]User{}, err
	}
	return &tasks, nil

}

//UpdateUser ...
func (u *User) UpdateUser(uid uint32) (*User, error) {
	DB = DB.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"name":       u.Name,
			"dob":        u.DOB,
			"updated_at": time.Now(),
		},
	)

	if DB.Error != nil {
		return &User{}, DB.Error
	}

	return u, nil
}

// DeleteUser ...
func (u *User) DeleteUser(uid uint32) (int64, error) {
	DB = DB.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if DB.Error != nil {
		return 0, DB.Error
	}
	return DB.RowsAffected, nil
}

// SignIn ...
func (u *User) SignIn(email, password string) (map[string]interface{}, error) {

	var err error

	userData := make(map[string]interface{})

	user := User{}

	err = DB.Debug().Model(User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		fmt.Println("this is the error getting the user: ", err)
		return nil, err
	}
	err = security.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println("this is the error hashing the password: ", err)
		return nil, err
	}
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		fmt.Println("this is the error creating the token: ", err)
		return nil, err
	}
	userData["token"] = token
	userData["id"] = user.ID
	userData["email"] = user.Email

	return userData, nil
}
