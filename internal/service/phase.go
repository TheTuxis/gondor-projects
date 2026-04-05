package service

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/TheTuxis/gondor-projects/internal/model"
	"github.com/TheTuxis/gondor-projects/internal/repository"
)

var (
	ErrPhaseNotFound = errors.New("phase not found")
)

type PhaseService struct {
	phaseRepo *repository.PhaseRepository
	logger    *zap.Logger
}

func NewPhaseService(phaseRepo *repository.PhaseRepository, logger *zap.Logger) *PhaseService {
	return &PhaseService{phaseRepo: phaseRepo, logger: logger}
}

func (s *PhaseService) List(projectID uint) ([]model.Phase, error) {
	return s.phaseRepo.ListByProject(projectID)
}

func (s *PhaseService) Create(projectID uint, input model.PhaseCreate) (*model.Phase, error) {
	status := input.Status
	if status == "" {
		status = "pending"
	}

	phase := &model.Phase{
		ProjectID:   projectID,
		Name:        input.Name,
		Description: input.Description,
		Order:       input.Order,
		Status:      status,
	}

	if input.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *input.StartDate); err == nil {
			phase.StartDate = &t
		}
	}
	if input.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *input.EndDate); err == nil {
			phase.EndDate = &t
		}
	}

	if err := s.phaseRepo.Create(phase); err != nil {
		return nil, err
	}

	return phase, nil
}

func (s *PhaseService) Update(id uint, input model.PhaseUpdate) (*model.Phase, error) {
	phase, err := s.phaseRepo.FindByID(id)
	if err != nil {
		return nil, ErrPhaseNotFound
	}

	if input.Name != nil {
		phase.Name = *input.Name
	}
	if input.Description != nil {
		phase.Description = *input.Description
	}
	if input.Order != nil {
		phase.Order = *input.Order
	}
	if input.Status != nil {
		phase.Status = *input.Status
	}
	if input.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *input.StartDate); err == nil {
			phase.StartDate = &t
		}
	}
	if input.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *input.EndDate); err == nil {
			phase.EndDate = &t
		}
	}

	if err := s.phaseRepo.Update(phase); err != nil {
		return nil, err
	}

	return phase, nil
}

func (s *PhaseService) Delete(id uint) error {
	if _, err := s.phaseRepo.FindByID(id); err != nil {
		return ErrPhaseNotFound
	}
	return s.phaseRepo.Delete(id)
}
