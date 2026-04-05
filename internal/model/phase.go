package model

import (
	"time"
)

type Phase struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ProjectID   uint       `gorm:"index;not null" json:"project_id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `json:"description"`
	Order       int        `gorm:"default:0" json:"order"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      string     `gorm:"default:pending" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PhaseCreate struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Order       int     `json:"order"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Status      string  `json:"status"`
}

type PhaseUpdate struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Order       *int    `json:"order"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Status      *string `json:"status"`
}
