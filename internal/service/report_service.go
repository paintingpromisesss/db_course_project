package service

import (
	"context"

	"db_course_project/internal/models"
	"db_course_project/internal/pagination"
	"db_course_project/internal/repository"
)

type ReportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) ActiveRosters(ctx context.Context, limit, offset int) ([]models.ActiveRosterView, int, error) {
	limit, offset = pagination.Normalize(limit, offset)
	return s.repo.ActiveRosters(ctx, limit, offset)
}

func (s *ReportService) MatchResults(ctx context.Context, tournamentID *int64, limit, offset int) ([]models.MatchResultView, int, error) {
	limit, offset = pagination.Normalize(limit, offset)
	return s.repo.MatchResults(ctx, tournamentID, limit, offset)
}

func (s *ReportService) PlayerCareer(ctx context.Context, search string, limit, offset int) ([]models.PlayerCareerStats, int, error) {
	limit, offset = pagination.Normalize(limit, offset)
	return s.repo.PlayerCareer(ctx, search, limit, offset)
}

func (s *ReportService) TournamentStandings(ctx context.Context, tournamentID int64) ([]models.TournamentStanding, error) {
	return s.repo.TournamentStandings(ctx, tournamentID)
}

func (s *ReportService) PlayerKDA(ctx context.Context, playerID int64) (float64, error) {
	return s.repo.PlayerKDA(ctx, playerID)
}
