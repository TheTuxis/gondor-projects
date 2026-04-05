package repository

import (
	"github.com/TheTuxis/gondor-projects/internal/model"
	"gorm.io/gorm"
)

type MemberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) ListByProject(projectID uint) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	err := r.db.Where("project_id = ?", projectID).Find(&members).Error
	return members, err
}

func (r *MemberRepository) Create(member *model.ProjectMember) error {
	return r.db.Create(member).Error
}

func (r *MemberRepository) Delete(id uint) error {
	return r.db.Delete(&model.ProjectMember{}, id).Error
}

func (r *MemberRepository) FindByID(id uint) (*model.ProjectMember, error) {
	var member model.ProjectMember
	if err := r.db.First(&member, id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}
