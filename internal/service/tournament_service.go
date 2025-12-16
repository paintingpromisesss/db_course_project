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

type TournamentService struct {
	repo repository.TournamentRepository
}

func NewTournamentService(repo repository.TournamentRepository) *TournamentService {
	return &TournamentService{repo: repo}
}

func (s *TournamentService) Create(ctx context.Context, t *models.Tournament) error {
	t.Name = strings.TrimSpace(t.Name)
	t.Currency = strings.TrimSpace(t.Currency)
	t.Status = strings.TrimSpace(t.Status)
	if t.Name == "" || t.DisciplineID == 0 {
		return errors.New("name and discipline_id are required")
	}
	if t.EndDate.Before(t.StartDate) {
		return errors.New("end_date must be after start_date")
	}
	if t.StartDate.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return errors.New("start_date looks invalid")
	}
	return s.repo.Create(ctx, t)
}

func (s *TournamentService) Get(ctx context.Context, id int64) (*models.Tournament, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TournamentService) List(ctx context.Context, filter models.TournamentFilter) ([]models.Tournament, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *TournamentService) Update(ctx context.Context, t *models.Tournament) error {
	t.Name = strings.TrimSpace(t.Name)
	t.Currency = strings.TrimSpace(t.Currency)
	t.Status = strings.TrimSpace(t.Status)
	if t.Name == "" || t.DisciplineID == 0 {
		return errors.New("name and discipline_id are required")
	}
	if t.EndDate.Before(t.StartDate) {
		return errors.New("end_date must be after start_date")
	}
	return s.repo.Update(ctx, t)
}

func (s *TournamentService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
