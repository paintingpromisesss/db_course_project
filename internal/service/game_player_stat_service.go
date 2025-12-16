package service

import (
	"context"
	"errors"

	"db_course_project/internal/models"
	"db_course_project/internal/pagination"
	"db_course_project/internal/repository"
)

type GamePlayerStatService struct {
	repo repository.GamePlayerStatRepository
}

func NewGamePlayerStatService(repo repository.GamePlayerStatRepository) *GamePlayerStatService {
	return &GamePlayerStatService{repo: repo}
}

func (s *GamePlayerStatService) Create(ctx context.Context, st *models.GamePlayerStat) error {
	if st.GameID == 0 || st.PlayerID == 0 {
		return errors.New("game_id and player_id are required")
	}
	return s.repo.Create(ctx, st)
}

func (s *GamePlayerStatService) Get(ctx context.Context, id int64) (*models.GamePlayerStat, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *GamePlayerStatService) List(ctx context.Context, filter models.GamePlayerStatFilter) ([]models.GamePlayerStat, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *GamePlayerStatService) Update(ctx context.Context, st *models.GamePlayerStat) error {
	if st.GameID == 0 || st.PlayerID == 0 {
		return errors.New("game_id and player_id are required")
	}
	return s.repo.Update(ctx, st)
}

func (s *GamePlayerStatService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
