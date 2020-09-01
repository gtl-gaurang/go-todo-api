package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Task ... Task table structure
type Task struct {
	gorm.Model
	ID          uint32 `gorm:"primary_key;auto_increment" json:"id"`
	UserID      int    `gorm:"not null" json:"user_id"`
	User        User
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"size:255;not null" json:"description"`
	IsCompleted bool      `gorm:"size:1; default:0" json:"is_completed"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//Prepare ... Take data double check
func (t *Task) Prepare() {
	t.Name = html.EscapeString(strings.TrimSpace(t.Name))
	t.Description = html.EscapeString(strings.TrimSpace(t.Description))
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

// Validate ... Take data
func (t *Task) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "add":
		if t.Name == "" {
			err = errors.New("Required Task Name")
			errorMessages["required_name"] = err.Error()
		}
	default:
		if t.Name == "" {
			err = errors.New("Required Task Name")
			errorMessages["required_name"] = err.Error()
		}
	}

	return errorMessages
}

//AddTask ... Add task into DB
func (t *Task) AddTask(db *gorm.DB) (*Task, error) {
	var err error
	err = db.Debug().Create(&t).Error
	if err != nil {
		return &Task{}, err
	}
	return t, nil
}

//GetAllTask ... Get All task
func (t *Task) GetAllTask(db *gorm.DB) (*[]Task, error) {
	var err error
	tasks := []Task{}
	err = db.Debug().Model(&Task{}).Limit(100).Find(&tasks).Error
	if err != nil {
		return &[]Task{}, err
	}
	return &tasks, nil

}

//UpdateTask ...
func (t *Task) UpdateTask(db *gorm.DB, uid uint32) (*Task, error) {
	db = db.Debug().Model(&Task{}).Where("id = ?", uid).Take(&Task{}).UpdateColumns(
		map[string]interface{}{
			"description":  t.Description,
			"name":         t.Name,
			"is_completed": t.IsCompleted,
			"updated_at":   time.Now(),
		},
	)

	if db.Error != nil {
		return &Task{}, db.Error
	}

	return t, nil
}

// DeleteTask ...
func (t *Task) DeleteTask(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&Task{}).Where("id = ?", uid).Take(&Task{}).Delete(&Task{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
