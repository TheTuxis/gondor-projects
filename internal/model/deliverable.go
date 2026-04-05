package model

import (
	"time"
)

type Deliverable struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"index;not null" json:"project_id"`
	TaskID      *uint     `gorm:"index" json:"task_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Status      string    `gorm:"default:pending" json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeliverableCreate struct {
	TaskID      *uint   `json:"task_id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	DueDate     *string `json:"due_date"`
}

type DeliverableUpdate struct {
	TaskID      *uint   `json:"task_id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	DueDate     *string `json:"due_date"`
}
