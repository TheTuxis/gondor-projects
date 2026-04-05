package repository

import (
	"github.com/TheTuxis/gondor-projects/internal/model"
	"gorm.io/gorm"
)

type DeliverableRepository struct {
	db *gorm.DB
}

func NewDeliverableRepository(db *gorm.DB) *DeliverableRepository {
	return &DeliverableRepository{db: db}
}

func (r *DeliverableRepository) FindByID(id uint) (*model.Deliverable, error) {
	var deliverable model.Deliverable
	if err := r.db.First(&deliverable, id).Error; err != nil {
		return nil, err
	}
	return &deliverable, nil
}

func (r *DeliverableRepository) ListByProject(projectID uint) ([]model.Deliverable, error) {
	var deliverables []model.Deliverable
	err := r.db.Where("project_id = ?", projectID).Find(&deliverables).Error
	return deliverables, err
}

func (r *DeliverableRepository) Create(deliverable *model.Deliverable) error {
	return r.db.Create(deliverable).Error
}

func (r *DeliverableRepository) Update(deliverable *model.Deliverable) error {
	return r.db.Save(deliverable).Error
}

func (r *DeliverableRepository) Delete(id uint) error {
	return r.db.Delete(&model.Deliverable{}, id).Error
}
