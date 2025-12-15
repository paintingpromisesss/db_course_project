package service

import (
	"context"
	"errors"
	"strings"

	"db_course_project/internal/models"
	"db_course_project/internal/pagination"
	"db_course_project/internal/repository"
)

// DisciplineService orchestrates discipline use cases.
type DisciplineService struct {
	repo repository.DisciplineRepository
}

func NewDisciplineService(repo repository.DisciplineRepository) *DisciplineService {
	return &DisciplineService{repo: repo}
}

func (s *DisciplineService) Create(ctx context.Context, d *models.Discipline) error {
	d.Code = strings.TrimSpace(d.Code)
	d.Name = strings.TrimSpace(d.Name)
	d.Description = strings.TrimSpace(d.Description)
	if d.Code == "" || d.Name == "" {
		return errors.New("code and name are required")
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return err
	}
	return nil
}

func (s *DisciplineService) Get(ctx context.Context, id int64) (*models.Discipline, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DisciplineService) List(ctx context.Context, filter models.DisciplineFilter) ([]models.Discipline, int, error) {
	filter.Limit, filter.Offset = pagination.Normalize(filter.Limit, filter.Offset)
	return s.repo.List(ctx, filter)
}

func (s *DisciplineService) Update(ctx context.Context, d *models.Discipline) error {
	d.Code = strings.TrimSpace(d.Code)
	d.Name = strings.TrimSpace(d.Name)
	d.Description = strings.TrimSpace(d.Description)
	if d.Code == "" || d.Name == "" {
		return errors.New("code and name are required")
	}
	return s.repo.Update(ctx, d)
}

func (s *DisciplineService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
