package service

import (
	"context"
	"errors"
	"strings"

	"db_course_project/internal/models"
	"db_course_project/internal/pagination"
	"db_course_project/internal/repository"
)

// TournamentRegistrationService orchestrates registrations.
type TournamentRegistrationService struct {
	repo repository.TournamentRegistrationRepository
}

func NewTournamentRegistrationService(repo repository.TournamentRegistrationRepository) *TournamentRegistrationService {
	return &TournamentRegistrationService{repo: repo}
}

func (s *TournamentRegistrationService) Create(ctx context.Context, reg *models.TournamentRegistration) error {
	reg.Status = strings.TrimSpace(reg.Status)
	if reg.Status == "" {
		reg.Status = "Pending"
	}
	if reg.TournamentID == 0 || reg.TeamID == 0 {
		return errors.New("tournament_id and team_id are required")
	}
	return s.repo.Create(ctx, reg)
}

func (s *TournamentRegistrationService) Get(ctx context.Context, id int64) (*models.TournamentRegistration, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TournamentRegistrationService) List(ctx context.Context, filter models.TournamentRegistrationFilter) ([]models.TournamentRegistration, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *TournamentRegistrationService) Update(ctx context.Context, reg *models.TournamentRegistration) error {
	reg.Status = strings.TrimSpace(reg.Status)
	if reg.Status == "" {
		reg.Status = "Pending"
	}
	if reg.TournamentID == 0 || reg.TeamID == 0 {
		return errors.New("tournament_id and team_id are required")
	}
	return s.repo.Update(ctx, reg)
}

func (s *TournamentRegistrationService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
