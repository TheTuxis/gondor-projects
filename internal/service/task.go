package service

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/TheTuxis/gondor-projects/internal/model"
	"github.com/TheTuxis/gondor-projects/internal/repository"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService struct {
	taskRepo *repository.TaskRepository
	logger   *zap.Logger
}

func NewTaskService(taskRepo *repository.TaskRepository, logger *zap.Logger) *TaskService {
	return &TaskService{taskRepo: taskRepo, logger: logger}
}

func (s *TaskService) List(projectID uint, params model.ListParams) (*model.PaginatedResult, error) {
	tasks, total, err := s.taskRepo.ListByProject(projectID, params)
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
		Data: tasks,
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

func (s *TaskService) GetByID(id uint) (*model.Task, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (s *TaskService) Create(projectID uint, input model.TaskCreate) (*model.Task, error) {
	status := input.Status
	if status == "" {
		status = "pending"
	}
	priority := input.Priority
	if priority == "" {
		priority = "medium"
	}

	task := &model.Task{
		ProjectID:      projectID,
		Name:           input.Name,
		Description:    input.Description,
		AssigneeID:     input.AssigneeID,
		Status:         status,
		Priority:       priority,
		EstimatedHours: input.EstimatedHours,
		ParentID:       input.ParentID,
	}

	if input.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *input.StartDate); err == nil {
			task.StartDate = &t
		}
	}
	if input.DueDate != nil {
		if t, err := time.Parse("2006-01-02", *input.DueDate); err == nil {
			task.DueDate = &t
		}
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) Update(id uint, input model.TaskUpdate) (*model.Task, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if input.Name != nil {
		task.Name = *input.Name
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.AssigneeID != nil {
		task.AssigneeID = input.AssigneeID
	}
	if input.Status != nil {
		task.Status = *input.Status
	}
	if input.Priority != nil {
		task.Priority = *input.Priority
	}
	if input.EstimatedHours != nil {
		task.EstimatedHours = input.EstimatedHours
	}
	if input.ActualHours != nil {
		task.ActualHours = input.ActualHours
	}
	if input.ParentID != nil {
		task.ParentID = input.ParentID
	}
	if input.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *input.StartDate); err == nil {
			task.StartDate = &t
		}
	}
	if input.DueDate != nil {
		if t, err := time.Parse("2006-01-02", *input.DueDate); err == nil {
			task.DueDate = &t
		}
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) Delete(id uint) error {
	if _, err := s.taskRepo.FindByID(id); err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.Delete(id)
}
