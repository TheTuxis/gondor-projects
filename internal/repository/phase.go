package repository

import (
	"github.com/TheTuxis/gondor-projects/internal/model"
	"gorm.io/gorm"
)

type PhaseRepository struct {
	db *gorm.DB
}

func NewPhaseRepository(db *gorm.DB) *PhaseRepository {
	return &PhaseRepository{db: db}
}

func (r *PhaseRepository) FindByID(id uint) (*model.Phase, error) {
	var phase model.Phase
	if err := r.db.First(&phase, id).Error; err != nil {
		return nil, err
	}
	return &phase, nil
}

func (r *PhaseRepository) ListByProject(projectID uint) ([]model.Phase, error) {
	var phases []model.Phase
	err := r.db.Where("project_id = ?", projectID).Order("\"order\" asc").Find(&phases).Error
	return phases, err
}

func (r *PhaseRepository) Create(phase *model.Phase) error {
	return r.db.Create(phase).Error
}

func (r *PhaseRepository) Update(phase *model.Phase) error {
	return r.db.Save(phase).Error
}

func (r *PhaseRepository) Delete(id uint) error {
	return r.db.Delete(&model.Phase{}, id).Error
}
