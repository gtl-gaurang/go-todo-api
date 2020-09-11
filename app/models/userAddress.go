package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"
)

// UserAddress ...
type UserAddress struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserID       uint32    `gorm:"not null" json:"user_id"`
	Title        string    `gorm:"size:255;not null" json:"title"`
	AddressLine1 string    `gorm:"size:255;null" json:"address_line1"`
	AddressLine2 string    `gorm:"size:255;null" json:"address_line2"`
	Country      string    `gorm:"size:100;null" json:"country"`
	State        string    `gorm:"size:100;null" json:"state"`
	City         string    `gorm:"size:100;null" json:"city"`
	Pin          string    `gorm:"size:10;null" json:"pin"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//Prepare ...
func (ua *UserAddress) Prepare() {
	ua.Title = html.EscapeString(strings.TrimSpace(ua.Title))
	ua.AddressLine1 = html.EscapeString(strings.TrimSpace(ua.AddressLine1))
	ua.AddressLine2 = html.EscapeString(strings.TrimSpace(ua.AddressLine2))
	ua.Country = html.EscapeString(strings.TrimSpace(ua.Country))
	ua.State = html.EscapeString(strings.TrimSpace(ua.State))
	ua.City = html.EscapeString(strings.TrimSpace(ua.City))
	ua.CreatedAt = time.Now()
	ua.UpdatedAt = time.Now()
}

// Validate ...
func (ua *UserAddress) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	default:
		if ua.Title == "" {
			err = errors.New("Required address title")
			errorMessages["required_title"] = err.Error()
		}
		if ua.AddressLine1 == "" {
			err = errors.New("Required address address Line1")
			errorMessages["required_address_line1"] = err.Error()
		}
		if ua.Country == "" {
			err = errors.New("Required country")
			errorMessages["required_country"] = err.Error()
		}
		if ua.State == "" {
			err = errors.New("Required state")
			errorMessages["required_state"] = err.Error()
		}
		if ua.City == "" {
			err = errors.New("Required city")
			errorMessages["required_city"] = err.Error()
		}
	}

	return errorMessages
}

//AddAddress ...
func (db *DataSource) AddAddress(ua *UserAddress) (*UserAddress, error) {
	var err error
	fmt.Println()
	err = db.DB.Debug().Create(&ua).Error
	if err != nil {
		return &UserAddress{}, err
	}
	return ua, nil
}

//GetAllAddress ...
func (db *DataSource) GetAllAddress(ua *UserAddress) (*[]UserAddress, error) {
	var err error
	address := []UserAddress{}
	err = db.DB.Debug().Model(&UserAddress{}).Limit(100).Find(&address).Error
	if err != nil {
		return &[]UserAddress{}, err
	}
	return &address, nil

}

//UpdateAddress ...
func (db *DataSource) UpdateAddress(aid uint32, ua *UserAddress) (*UserAddress, error) {
	db.DB = db.DB.Debug().Model(&UserAddress{}).Where("id = ?", aid).Take(&UserAddress{}).UpdateColumns(
		map[string]interface{}{
			"title":         ua.Title,
			"address_line1": ua.AddressLine1,
			"address_line2": ua.AddressLine2,
			"country":       ua.Country,
			"state":         ua.State,
			"city":          ua.City,
			"pin":           ua.Pin,
			"updated_at":    time.Now(),
		},
	)

	if db.DB.Error != nil {
		return &UserAddress{}, db.DB.Error
	}

	return ua, nil
}

// DeleteAddress ...
func (db *DataSource) DeleteAddress(aid uint32, ua *UserAddress) (int64, error) {
	db.DB = db.DB.Debug().Model(&UserAddress{}).Where("id = ?", aid).Take(&UserAddress{}).Delete(&User{})

	if db.DB.Error != nil {
		return 0, db.DB.Error
	}
	return db.DB.RowsAffected, nil
}
