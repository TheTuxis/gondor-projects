package service

import (
	"errors"

	"go.uber.org/zap"

	"github.com/TheTuxis/gondor-projects/internal/model"
	"github.com/TheTuxis/gondor-projects/internal/repository"
)

var (
	ErrMemberNotFound = errors.New("member not found")
)

type MemberService struct {
	memberRepo *repository.MemberRepository
	logger     *zap.Logger
}

func NewMemberService(memberRepo *repository.MemberRepository, logger *zap.Logger) *MemberService {
	return &MemberService{memberRepo: memberRepo, logger: logger}
}

func (s *MemberService) List(projectID uint) ([]model.ProjectMember, error) {
	return s.memberRepo.ListByProject(projectID)
}

func (s *MemberService) Create(projectID uint, input model.MemberCreate) (*model.ProjectMember, error) {
	role := input.Role
	if role == "" {
		role = "team_member"
	}

	member := &model.ProjectMember{
		ProjectID: projectID,
		UserID:    input.UserID,
		Role:      role,
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *MemberService) Delete(id uint) error {
	if _, err := s.memberRepo.FindByID(id); err != nil {
		return ErrMemberNotFound
	}
	return s.memberRepo.Delete(id)
}
