package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ProjectID      uint           `gorm:"index;not null" json:"project_id"`
	Name           string         `gorm:"not null" json:"name"`
	Description    string         `json:"description"`
	AssigneeID     *uint          `gorm:"index" json:"assignee_id"`
	Status         string         `gorm:"default:pending;index" json:"status"`
	Priority       string         `gorm:"default:medium" json:"priority"`
	StartDate      *time.Time     `json:"start_date"`
	DueDate        *time.Time     `json:"due_date"`
	EstimatedHours *float64       `json:"estimated_hours"`
	ActualHours    *float64       `json:"actual_hours"`
	ParentID       *uint          `gorm:"index" json:"parent_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Subtasks []Task `gorm:"foreignKey:ParentID" json:"subtasks,omitempty"`
}

type TaskCreate struct {
	Name           string   `json:"name" binding:"required"`
	Description    string   `json:"description"`
	AssigneeID     *uint    `json:"assignee_id"`
	Status         string   `json:"status"`
	Priority       string   `json:"priority"`
	StartDate      *string  `json:"start_date"`
	DueDate        *string  `json:"due_date"`
	EstimatedHours *float64 `json:"estimated_hours"`
	ParentID       *uint    `json:"parent_id"`
}

type TaskUpdate struct {
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	AssigneeID     *uint    `json:"assignee_id"`
	Status         *string  `json:"status"`
	Priority       *string  `json:"priority"`
	StartDate      *string  `json:"start_date"`
	DueDate        *string  `json:"due_date"`
	EstimatedHours *float64 `json:"estimated_hours"`
	ActualHours    *float64 `json:"actual_hours"`
	ParentID       *uint    `json:"parent_id"`
}
