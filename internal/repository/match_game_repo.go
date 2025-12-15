package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"

	"db_course_project/internal/models"
)

// MatchGameRepository persistence.
type MatchGameRepository interface {
	Create(ctx context.Context, g *models.MatchGame) error
	GetByID(ctx context.Context, id int64) (*models.MatchGame, error)
	List(ctx context.Context, filter models.MatchGameFilter) ([]models.MatchGame, int, error)
	Update(ctx context.Context, g *models.MatchGame) error
	Delete(ctx context.Context, id int64) error
}

func NewMatchGameRepository(db *sqlx.DB) MatchGameRepository {
	return &matchGameRepo{db: db}
}

// ErrMatchGameNotFound signals missing row.
var ErrMatchGameNotFound = errors.New("match game not found")

type matchGameRepo struct {
	db *sqlx.DB
}

func (r *matchGameRepo) Create(ctx context.Context, g *models.MatchGame) error {
	query := `INSERT INTO match_games (match_id, map_name, game_number, duration_seconds, winner_team_id, score_team1, score_team2, started_at, had_technical_pause, pick_ban_phase)
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`
	return r.db.QueryRowxContext(ctx, query,
		g.MatchID,
		g.MapName,
		g.GameNumber,
		g.DurationSeconds,
		g.WinnerTeamID,
		g.ScoreTeam1,
		g.ScoreTeam2,
		g.StartedAt,
		g.HadTechnicalPause,
		g.PickBanPhase,
	).Scan(&g.ID)
}

func (r *matchGameRepo) GetByID(ctx context.Context, id int64) (*models.MatchGame, error) {
	var g models.MatchGame
	query := `SELECT id, match_id, map_name, game_number, duration_seconds, winner_team_id, score_team1, score_team2, started_at, had_technical_pause, pick_ban_phase
			  FROM match_games WHERE id=$1`
	if err := r.db.GetContext(ctx, &g, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMatchGameNotFound
		}
		return nil, err
	}
	return &g, nil
}

func (r *matchGameRepo) List(ctx context.Context, filter models.MatchGameFilter) ([]models.MatchGame, int, error) {
	base := `FROM match_games WHERE 1=1`
	args := []any{}
	conds := strings.Builder{}

	if filter.MatchID != nil {
		args = append(args, *filter.MatchID)
		conds.WriteString(` AND match_id = $` + strconv.Itoa(len(args)))
	}
	if filter.WinnerTeamID != nil {
		args = append(args, *filter.WinnerTeamID)
		conds.WriteString(` AND winner_team_id = $` + strconv.Itoa(len(args)))
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT id, match_id, map_name, game_number, duration_seconds, winner_team_id, score_team1, score_team2, started_at, had_technical_pause, pick_ban_phase ` + base + conds.String() +
		` ORDER BY match_id DESC, game_number ASC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.MatchGame{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *matchGameRepo) Update(ctx context.Context, g *models.MatchGame) error {
	query := `UPDATE match_games SET match_id=$1, map_name=$2, game_number=$3, duration_seconds=$4, winner_team_id=$5, score_team1=$6, score_team2=$7, started_at=$8, had_technical_pause=$9, pick_ban_phase=$10
			  WHERE id=$11`
	res, err := r.db.ExecContext(ctx, query,
		g.MatchID,
		g.MapName,
		g.GameNumber,
		g.DurationSeconds,
		g.WinnerTeamID,
		g.ScoreTeam1,
		g.ScoreTeam2,
		g.StartedAt,
		g.HadTechnicalPause,
		g.PickBanPhase,
		g.ID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrMatchGameNotFound
	}
	return nil
}

func (r *matchGameRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM match_games WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrMatchGameNotFound
	}
	return nil
}
