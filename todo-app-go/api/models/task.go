package models

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	Name        string    `gorm:"not null;type:varchar(100)" json:"name"`
	Status      string    `gorm:"not null;type:varchar(100)" json:"status"`
	Description string    `gorm:"type:varchar(500)" json:"description"`
	DueDate     time.Time `gorm:"not null;type:date" json:"dueDate"`
	// CreatedBy   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	UserID uint `gorm:"not null"                 json:"userId"`
}

func (task *Task) CreateTask(db *gorm.DB) error {
	return db.Debug().Create(&task).Error
}

func GetTasksByUserId(userId uint, db *gorm.DB) *[]Task {
	tasksList := []Task{}
	db.Debug().Where("user_id = ?", userId).Find(&tasksList)
	return &tasksList
}

func GetTasksByTaskId(taskId int, db *gorm.DB) *Task {
	recordFound := Task{}
	db.Debug().Where("id = ?", taskId).Find(&recordFound)
	return &recordFound
}

func (task *Task) DeleteTask(db *gorm.DB) error {
	return db.Debug().Where("id = ?", task.ID).Unscoped().Delete(&task).Error
}

func (task *Task) UpdateTask(taskId int, db *gorm.DB) error {
	updatedTask := Task{Name: task.Name, Status: task.Status, Description: task.Description, DueDate: task.DueDate}
	return db.Debug().Table("tasks").Where("id = ?", taskId).Updates(&updatedTask).Error
}

func (task *Task) PrepareData() {
	task.Name = strings.TrimSpace(task.Name)
	task.Description = strings.TrimSpace(task.Description)
	task.Status = strings.TrimSpace(task.Status)
}

func (task *Task) ValidateFields() error {
	if task.Name == "" {
		return errors.New("task name is required")
	}
	if task.Description == "" {
		return errors.New("task description is required")
	}
	if task.Status == "" {
		return errors.New("task status is required")
	}
	if task.UserID == 0 {
		return errors.New("task assignee id is required")
	}
	return nil
}
