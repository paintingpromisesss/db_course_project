package service

import (
	"context"
	"errors"
	"strings"

	"db_course_project/internal/models"
	"db_course_project/internal/pagination"
	"db_course_project/internal/repository"
)

// MatchGameService orchestrates per-map data.
type MatchGameService struct {
	repo repository.MatchGameRepository
}

func NewMatchGameService(repo repository.MatchGameRepository) *MatchGameService {
	return &MatchGameService{repo: repo}
}

func (s *MatchGameService) Create(ctx context.Context, g *models.MatchGame) error {
	g.MapName = strings.TrimSpace(g.MapName)
	if g.MatchID == 0 || g.MapName == "" || g.GameNumber <= 0 {
		return errors.New("match_id, map_name, game_number are required")
	}
	return s.repo.Create(ctx, g)
}

func (s *MatchGameService) Get(ctx context.Context, id int64) (*models.MatchGame, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MatchGameService) List(ctx context.Context, filter models.MatchGameFilter) ([]models.MatchGame, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *MatchGameService) Update(ctx context.Context, g *models.MatchGame) error {
	g.MapName = strings.TrimSpace(g.MapName)
	if g.MatchID == 0 || g.MapName == "" || g.GameNumber <= 0 {
		return errors.New("match_id, map_name, game_number are required")
	}
	return s.repo.Update(ctx, g)
}

func (s *MatchGameService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
