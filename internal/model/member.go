package model

import (
	"time"
)

type ProjectMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"index;not null;uniqueIndex:idx_project_user" json:"project_id"`
	UserID    uint      `gorm:"index;not null;uniqueIndex:idx_project_user" json:"user_id"`
	Role      string    `gorm:"default:team_member" json:"role"`
	JoinedAt  time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

type MemberCreate struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role"`
}
