package service

import (
	"context"
	"errors"
	"strings"

	"db_course_project/internal/models"
	"db_course_project/internal/pagination"
	"db_course_project/internal/repository"
)

// TeamService orchestrates team use cases.
type TeamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) Create(ctx context.Context, t *models.Team) error {
	t.Name = strings.TrimSpace(t.Name)
	t.Tag = strings.TrimSpace(t.Tag)
	t.CountryCode = strings.TrimSpace(t.CountryCode)
	if t.Name == "" || t.Tag == "" || t.CountryCode == "" || t.DisciplineID == 0 {
		return errors.New("name, tag, country_code, discipline_id are required")
	}
	return s.repo.Create(ctx, t)
}

func (s *TeamService) Get(ctx context.Context, id int64) (*models.Team, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TeamService) List(ctx context.Context, filter models.TeamFilter) ([]models.Team, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *TeamService) Update(ctx context.Context, t *models.Team) error {
	t.Name = strings.TrimSpace(t.Name)
	t.Tag = strings.TrimSpace(t.Tag)
	t.CountryCode = strings.TrimSpace(t.CountryCode)
	if t.Name == "" || t.Tag == "" || t.CountryCode == "" || t.DisciplineID == 0 {
		return errors.New("name, tag, country_code, discipline_id are required")
	}
	return s.repo.Update(ctx, t)
}

func (s *TeamService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
