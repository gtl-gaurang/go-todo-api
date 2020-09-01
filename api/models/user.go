package models

import (
	"errors"
	"html"
	"strings"
	"time"
	"todo-api/api/security"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
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
	default:
		if u.Name == "" {
			err = errors.New("Required User full Name")
			errorMessages["required_name"] = err.Error()
		}
	}

	return errorMessages
}

//AddUser ...
func (u *User) AddUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

//GetAllUser ... Get All user
func (u *User) GetAllUser(db *gorm.DB) (*[]User, error) {
	var err error
	tasks := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&tasks).Error
	if err != nil {
		return &[]User{}, err
	}
	return &tasks, nil

}

//UpdateUser ...
func (u *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"name":       u.Name,
			"dob":        u.DOB,
			"updated_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	return u, nil
}

// DeleteUser ...
func (u *User) DeleteUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
