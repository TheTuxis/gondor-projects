package model

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	CompanyID   uint           `gorm:"index;not null" json:"company_id"`
	ProjectType string         `gorm:"default:standard" json:"project_type"`
	Status      string         `gorm:"default:active;index" json:"status"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Budget      *float64       `json:"budget"`
	Currency    string         `gorm:"default:USD" json:"currency"`
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Tasks        []Task           `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
	Phases       []Phase          `gorm:"foreignKey:ProjectID" json:"phases,omitempty"`
	Members      []ProjectMember  `gorm:"foreignKey:ProjectID" json:"members,omitempty"`
	Deliverables []Deliverable    `gorm:"foreignKey:ProjectID" json:"deliverables,omitempty"`
}

type ProjectCreate struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	CompanyID   uint     `json:"company_id" binding:"required"`
	ProjectType string   `json:"project_type"`
	Status      string   `json:"status"`
	StartDate   *string  `json:"start_date"`
	EndDate     *string  `json:"end_date"`
	Budget      *float64 `json:"budget"`
	Currency    string   `json:"currency"`
}

type ProjectUpdate struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	ProjectType *string  `json:"project_type"`
	Status      *string  `json:"status"`
	StartDate   *string  `json:"start_date"`
	EndDate     *string  `json:"end_date"`
	Budget      *float64 `json:"budget"`
	Currency    *string  `json:"currency"`
}
