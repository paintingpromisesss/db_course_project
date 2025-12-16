package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"db_course_project/internal/models"
	"db_course_project/internal/pagination"
	"db_course_project/internal/repository"
)

type SquadMemberService struct {
	repo repository.SquadMemberRepository
}

func NewSquadMemberService(repo repository.SquadMemberRepository) *SquadMemberService {
	return &SquadMemberService{repo: repo}
}

func (s *SquadMemberService) validateDates(m *models.SquadMember) error {
	if m.LeaveDate != nil && m.LeaveDate.Before(m.JoinDate) {
		return errors.New("leave_date cannot be before join_date")
	}
	if m.ContractEndDate != nil && m.ContractEndDate.Before(m.JoinDate) {
		return errors.New("contract_end_date cannot be before join_date")
	}
	return nil
}

func (s *SquadMemberService) Create(ctx context.Context, m *models.SquadMember) error {
	m.Role = strings.TrimSpace(m.Role)
	if m.TeamID == 0 || m.PlayerID == 0 {
		return errors.New("team_id and player_id are required")
	}
	if m.Role == "" {
		m.Role = "Player"
	}
	if m.JoinDate.IsZero() {
		m.JoinDate = time.Now().UTC().Truncate(24 * time.Hour)
	}
	if err := s.validateDates(m); err != nil {
		return err
	}
	return s.repo.Create(ctx, m)
}

func (s *SquadMemberService) Get(ctx context.Context, id int64) (*models.SquadMember, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SquadMemberService) List(ctx context.Context, filter models.SquadMemberFilter) ([]models.SquadMember, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *SquadMemberService) Update(ctx context.Context, m *models.SquadMember) error {
	m.Role = strings.TrimSpace(m.Role)
	if m.TeamID == 0 || m.PlayerID == 0 {
		return errors.New("team_id and player_id are required")
	}
	if m.Role == "" {
		m.Role = "Player"
	}
	if m.JoinDate.IsZero() {
		return errors.New("join_date is required")
	}
	if err := s.validateDates(m); err != nil {
		return err
	}
	return s.repo.Update(ctx, m)
}

func (s *SquadMemberService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
