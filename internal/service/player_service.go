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

type PlayerService struct {
	repo repository.PlayerRepository
}

func NewPlayerService(repo repository.PlayerRepository) *PlayerService {
	return &PlayerService{repo: repo}
}

func (s *PlayerService) Create(ctx context.Context, p *models.Player) error {
	p.Nickname = strings.TrimSpace(p.Nickname)
	if p.Nickname == "" {
		return errors.New("nickname is required")
	}
	if p.BirthDate != nil && p.BirthDate.After(time.Now()) {
		return errors.New("birth_date cannot be in the future")
	}
	return s.repo.Create(ctx, p)
}

func (s *PlayerService) Get(ctx context.Context, id int64) (*models.Player, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PlayerService) List(ctx context.Context, filter models.PlayerFilter) ([]models.Player, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *PlayerService) Update(ctx context.Context, p *models.Player) error {
	p.Nickname = strings.TrimSpace(p.Nickname)
	if p.Nickname == "" {
		return errors.New("nickname is required")
	}
	if p.BirthDate != nil && p.BirthDate.After(time.Now()) {
		return errors.New("birth_date cannot be in the future")
	}
	return s.repo.Update(ctx, p)
}

func (s *PlayerService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
