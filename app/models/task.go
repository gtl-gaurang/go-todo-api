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
	ID          uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserID      uint32    `gorm:"not null" json:"user_id"`
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
func (db *DataSource) AddTask(t *Task) (*Task, error) {
	var err error
	err = db.DB.Debug().Create(&t).Error
	if err != nil {
		return &Task{}, err
	}
	return t, nil
}

//GetAllTask ... Get All task
func (db *DataSource) GetAllTask(t *Task) (*[]Task, error) {
	var err error
	tasks := []Task{}
	err = db.DB.Debug().Model(&Task{}).Limit(100).Find(&tasks).Error
	if err != nil {
		return &[]Task{}, err
	}
	return &tasks, nil

}

//UpdateTask ...
func (db *DataSource) UpdateTask(tid uint32, t *Task) (*Task, error) {
	db.DB = db.DB.Debug().Model(&Task{}).Where("id = ?", tid).Take(&Task{}).UpdateColumns(
		map[string]interface{}{
			"description":  t.Description,
			"name":         t.Name,
			"is_completed": t.IsCompleted,
			"updated_at":   time.Now(),
		},
	)

	if db.DB.Error != nil {
		return &Task{}, db.DB.Error
	}

	return t, nil
}

//FindTaskByID ...
func (db *DataSource) FindTaskByID(tid uint32, t *Task) (*Task, error) {
	var err error
	err = db.DB.Debug().Model(Task{}).Where("id = ?", tid).Take(&t).Error
	if err != nil {
		return &Task{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Task{}, errors.New("Task Not Found")
	}
	return t, err
}

// DeleteTask ...
func (db *DataSource) DeleteTask(tid uint32, t *Task) (int64, error) {
	db.DB = db.DB.Debug().Model(&Task{}).Where("id = ?", tid).Take(&Task{}).Delete(&Task{})

	if db.DB.Error != nil {
		return 0, db.DB.Error
	}
	return db.DB.RowsAffected, nil
}
