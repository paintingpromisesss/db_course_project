package service

import (
	"context"
	"errors"

	"db_course_project/internal/models"
	"db_course_project/internal/pagination"
	"db_course_project/internal/repository"
)

type TeamProfileService struct {
	repo repository.TeamProfileRepository
}

func NewTeamProfileService(repo repository.TeamProfileRepository) *TeamProfileService {
	return &TeamProfileService{repo: repo}
}

func (s *TeamProfileService) Create(ctx context.Context, p *models.TeamProfile) error {
	if p.TeamID == 0 {
		return errors.New("team_id is required")
	}
	return s.repo.Create(ctx, p)
}

func (s *TeamProfileService) Get(ctx context.Context, teamID int64) (*models.TeamProfile, error) {
	return s.repo.GetByTeamID(ctx, teamID)
}

func (s *TeamProfileService) List(ctx context.Context, filter models.TeamProfileFilter) ([]models.TeamProfile, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *TeamProfileService) Update(ctx context.Context, p *models.TeamProfile) error {
	if p.TeamID == 0 {
		return errors.New("team_id is required")
	}
	return s.repo.Update(ctx, p)
}

func (s *TeamProfileService) Delete(ctx context.Context, teamID int64) error {
	return s.repo.Delete(ctx, teamID)
}
