package service

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/TheTuxis/gondor-projects/internal/model"
	"github.com/TheTuxis/gondor-projects/internal/repository"
)

var (
	ErrProjectNotFound = errors.New("project not found")
)

type ProjectService struct {
	projectRepo *repository.ProjectRepository
	logger      *zap.Logger
}

func NewProjectService(projectRepo *repository.ProjectRepository, logger *zap.Logger) *ProjectService {
	return &ProjectService{projectRepo: projectRepo, logger: logger}
}

func (s *ProjectService) List(params model.ListParams) (*model.PaginatedResult, error) {
	projects, total, err := s.projectRepo.List(params)
	if err != nil {
		return nil, err
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &model.PaginatedResult{
		Data: projects,
		Pagination: model.Pagination{
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
			HasNext:    params.Page < totalPages,
			HasPrev:    params.Page > 1,
		},
	}, nil
}

func (s *ProjectService) GetByID(id uint) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

func (s *ProjectService) Create(input model.ProjectCreate, createdBy uint) (*model.Project, error) {
	status := input.Status
	if status == "" {
		status = "active"
	}
	projectType := input.ProjectType
	if projectType == "" {
		projectType = "standard"
	}
	currency := input.Currency
	if currency == "" {
		currency = "USD"
	}

	project := &model.Project{
		Name:        input.Name,
		Description: input.Description,
		CompanyID:   input.CompanyID,
		ProjectType: projectType,
		Status:      status,
		Budget:      input.Budget,
		Currency:    currency,
		CreatedBy:   createdBy,
	}

	if input.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *input.StartDate); err == nil {
			project.StartDate = &t
		}
	}
	if input.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *input.EndDate); err == nil {
			project.EndDate = &t
		}
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) Update(id uint, input model.ProjectUpdate) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, ErrProjectNotFound
	}

	if input.Name != nil {
		project.Name = *input.Name
	}
	if input.Description != nil {
		project.Description = *input.Description
	}
	if input.ProjectType != nil {
		project.ProjectType = *input.ProjectType
	}
	if input.Status != nil {
		project.Status = *input.Status
	}
	if input.Budget != nil {
		project.Budget = input.Budget
	}
	if input.Currency != nil {
		project.Currency = *input.Currency
	}
	if input.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *input.StartDate); err == nil {
			project.StartDate = &t
		}
	}
	if input.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *input.EndDate); err == nil {
			project.EndDate = &t
		}
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) Delete(id uint) error {
	if _, err := s.projectRepo.FindByID(id); err != nil {
		return ErrProjectNotFound
	}
	return s.projectRepo.Delete(id)
}
