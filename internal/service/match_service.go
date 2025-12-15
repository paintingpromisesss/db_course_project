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

// MatchService orchestrates match operations.
type MatchService struct {
	repo repository.MatchRepository
}

func NewMatchService(repo repository.MatchRepository) *MatchService {
	return &MatchService{repo: repo}
}

func (s *MatchService) Create(ctx context.Context, m *models.Match) error {
	m.Format = strings.TrimSpace(m.Format)
	if m.Format == "" {
		m.Format = "bo3"
	}
	if m.TournamentID == 0 || m.StartTime.IsZero() {
		return errors.New("tournament_id and start_time are required")
	}
	if m.Team1ID != nil && m.Team2ID != nil && *m.Team1ID == *m.Team2ID {
		return errors.New("team1_id and team2_id must differ")
	}
	if m.StartTime.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return errors.New("start_time looks invalid")
	}
	return s.repo.Create(ctx, m)
}

func (s *MatchService) Get(ctx context.Context, id int64) (*models.Match, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MatchService) List(ctx context.Context, filter models.MatchFilter) ([]models.Match, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *MatchService) Update(ctx context.Context, m *models.Match) error {
	m.Format = strings.TrimSpace(m.Format)
	if m.Format == "" {
		m.Format = "bo3"
	}
	if m.TournamentID == 0 || m.StartTime.IsZero() {
		return errors.New("tournament_id and start_time are required")
	}
	if m.Team1ID != nil && m.Team2ID != nil && *m.Team1ID == *m.Team2ID {
		return errors.New("team1_id and team2_id must differ")
	}
	return s.repo.Update(ctx, m)
}

func (s *MatchService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
