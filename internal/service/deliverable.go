package service

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/TheTuxis/gondor-projects/internal/model"
	"github.com/TheTuxis/gondor-projects/internal/repository"
)

var (
	ErrDeliverableNotFound = errors.New("deliverable not found")
)

type DeliverableService struct {
	deliverableRepo *repository.DeliverableRepository
	logger          *zap.Logger
}

func NewDeliverableService(deliverableRepo *repository.DeliverableRepository, logger *zap.Logger) *DeliverableService {
	return &DeliverableService{deliverableRepo: deliverableRepo, logger: logger}
}

func (s *DeliverableService) List(projectID uint) ([]model.Deliverable, error) {
	return s.deliverableRepo.ListByProject(projectID)
}

func (s *DeliverableService) Create(projectID uint, input model.DeliverableCreate) (*model.Deliverable, error) {
	status := input.Status
	if status == "" {
		status = "pending"
	}

	deliverable := &model.Deliverable{
		ProjectID:   projectID,
		TaskID:      input.TaskID,
		Name:        input.Name,
		Description: input.Description,
		Status:      status,
	}

	if input.DueDate != nil {
		if t, err := time.Parse("2006-01-02", *input.DueDate); err == nil {
			deliverable.DueDate = &t
		}
	}

	if err := s.deliverableRepo.Create(deliverable); err != nil {
		return nil, err
	}

	return deliverable, nil
}

func (s *DeliverableService) Update(id uint, input model.DeliverableUpdate) (*model.Deliverable, error) {
	deliverable, err := s.deliverableRepo.FindByID(id)
	if err != nil {
		return nil, ErrDeliverableNotFound
	}

	if input.TaskID != nil {
		deliverable.TaskID = input.TaskID
	}
	if input.Name != nil {
		deliverable.Name = *input.Name
	}
	if input.Description != nil {
		deliverable.Description = *input.Description
	}
	if input.Status != nil {
		deliverable.Status = *input.Status
	}
	if input.DueDate != nil {
		if t, err := time.Parse("2006-01-02", *input.DueDate); err == nil {
			deliverable.DueDate = &t
		}
	}

	if err := s.deliverableRepo.Update(deliverable); err != nil {
		return nil, err
	}

	return deliverable, nil
}

func (s *DeliverableService) Delete(id uint) error {
	if _, err := s.deliverableRepo.FindByID(id); err != nil {
		return ErrDeliverableNotFound
	}
	return s.deliverableRepo.Delete(id)
}
