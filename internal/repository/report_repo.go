package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"

	"db_course_project/internal/models"
)

type ReportRepository interface {
	ActiveRosters(ctx context.Context, limit, offset int) ([]models.ActiveRosterView, int, error)
	MatchResults(ctx context.Context, tournamentID *int64, limit, offset int) ([]models.MatchResultView, int, error)
	PlayerCareer(ctx context.Context, search string, limit, offset int) ([]models.PlayerCareerStats, int, error)
	TournamentStandings(ctx context.Context, tournamentID int64) ([]models.TournamentStanding, error)
	PlayerKDA(ctx context.Context, playerID int64) (float64, error)
}

func NewReportRepository(db *sqlx.DB) ReportRepository {
	return &reportRepo{db: db}
}

type reportRepo struct {
	db *sqlx.DB
}

func (r *reportRepo) ActiveRosters(ctx context.Context, limit, offset int) ([]models.ActiveRosterView, int, error) {
	countQuery := `SELECT count(*) FROM v_active_rosters`
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	query := `SELECT team_id, team_name, tag, player_id, nickname, country_code, role, join_date
			   FROM v_active_rosters
			   ORDER BY team_name ASC, nickname ASC
			   LIMIT $1 OFFSET $2`
	rows := []models.ActiveRosterView{}
	if err := r.db.SelectContext(ctx, &rows, query, limit, offset); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *reportRepo) MatchResults(ctx context.Context, tournamentID *int64, limit, offset int) ([]models.MatchResultView, int, error) {
	base := `FROM v_match_results WHERE 1=1`
	args := []any{}
	conds := strings.Builder{}

	if tournamentID != nil {
		args = append(args, *tournamentID)
		conds.WriteString(` AND tournament_id = $` + strconv.Itoa(len(args)))
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT match_id, tournament_id, start_time, stage, format, winner_team_id, games_played, total_score_team1, total_score_team2 ` + base + conds.String() + `
				 ORDER BY start_time DESC, match_id DESC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.MatchResultView{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *reportRepo) PlayerCareer(ctx context.Context, search string, limit, offset int) ([]models.PlayerCareerStats, int, error) {
	base := `FROM v_player_career_stats WHERE 1=1`
	args := []any{}
	conds := strings.Builder{}
	if search != "" {
		args = append(args, "%"+search+"%")
		conds.WriteString(` AND LOWER(nickname) LIKE LOWER($` + strconv.Itoa(len(args)) + `)`)
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT player_id, nickname, kills, deaths, assists, damage, gold, kda ` + base + conds.String() + `
				 ORDER BY kda DESC, kills DESC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.PlayerCareerStats{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *reportRepo) TournamentStandings(ctx context.Context, tournamentID int64) ([]models.TournamentStanding, error) {
	query := `SELECT team_id, matches_played, wins, losses, forfeits FROM fn_tournament_standings($1)`
	rows := []models.TournamentStanding{}
	if err := r.db.SelectContext(ctx, &rows, query, tournamentID); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *reportRepo) PlayerKDA(ctx context.Context, playerID int64) (float64, error) {
	query := `SELECT fn_player_kda($1)`
	var kda float64
	if err := r.db.GetContext(ctx, &kda, query, playerID); err != nil {
		return 0, err
	}
	return kda, nil
}
