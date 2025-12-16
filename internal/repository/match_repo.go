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

// MatchRepository persistence.
type MatchRepository interface {
	Create(ctx context.Context, m *models.Match) error
	GetByID(ctx context.Context, id int64) (*models.Match, error)
	List(ctx context.Context, filter models.MatchFilter) ([]models.Match, int, error)
	Update(ctx context.Context, m *models.Match) error
	Delete(ctx context.Context, id int64) error
}

func NewMatchRepository(db *sqlx.DB) MatchRepository {
	return &matchRepo{db: db}
}

// ErrMatchNotFound signals missing row.
var ErrMatchNotFound = errors.New("match not found")

type matchRepo struct {
	db *sqlx.DB
}

func (r *matchRepo) Create(ctx context.Context, m *models.Match) error {
	query := `INSERT INTO matches (tournament_id, team1_id, team2_id, start_time, format, stage, winner_team_id, is_forfeit, match_notes)
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`
	return r.db.QueryRowxContext(ctx, query,
		m.TournamentID,
		m.Team1ID,
		m.Team2ID,
		m.StartTime,
		m.Format,
		m.Stage,
		m.WinnerTeamID,
		m.IsForfeit,
		m.MatchNotes,
	).Scan(&m.ID)
}

func (r *matchRepo) GetByID(ctx context.Context, id int64) (*models.Match, error) {
	var m models.Match
	query := `SELECT id, tournament_id, team1_id, team2_id, start_time, format, stage, winner_team_id, is_forfeit, match_notes
			  FROM matches WHERE id=$1`
	if err := r.db.GetContext(ctx, &m, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *matchRepo) List(ctx context.Context, filter models.MatchFilter) ([]models.Match, int, error) {
	base := `FROM matches WHERE 1=1`
	args := []any{}
	conds := strings.Builder{}

	if filter.TournamentID != nil {
		args = append(args, *filter.TournamentID)
		conds.WriteString(` AND tournament_id = $` + strconv.Itoa(len(args)))
	}
	if filter.TeamID != nil {
		args = append(args, *filter.TeamID, *filter.TeamID)
		conds.WriteString(` AND (team1_id = $` + strconv.Itoa(len(args)-1) + ` OR team2_id = $` + strconv.Itoa(len(args)) + `)`)
	}
	if filter.Stage != "" {
		args = append(args, filter.Stage)
		conds.WriteString(` AND LOWER(stage) = LOWER($` + strconv.Itoa(len(args)) + `)`)
	}
	if filter.Format != "" {
		args = append(args, filter.Format)
		conds.WriteString(` AND LOWER(format) = LOWER($` + strconv.Itoa(len(args)) + `)`)
	}
	if filter.From != nil {
		args = append(args, *filter.From)
		conds.WriteString(` AND start_time >= $` + strconv.Itoa(len(args)))
	}
	if filter.To != nil {
		args = append(args, *filter.To)
		conds.WriteString(` AND start_time <= $` + strconv.Itoa(len(args)))
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT id, tournament_id, team1_id, team2_id, start_time, format, stage, winner_team_id, is_forfeit, match_notes ` + base + conds.String() +
		` ORDER BY start_time DESC, id DESC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.Match{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *matchRepo) Update(ctx context.Context, m *models.Match) error {
	query := `UPDATE matches SET tournament_id=$1, team1_id=$2, team2_id=$3, start_time=$4, format=$5, stage=$6, winner_team_id=$7, is_forfeit=$8, match_notes=$9
			  WHERE id=$10`
	res, err := r.db.ExecContext(ctx, query,
		m.TournamentID,
		m.Team1ID,
		m.Team2ID,
		m.StartTime,
		m.Format,
		m.Stage,
		m.WinnerTeamID,
		m.IsForfeit,
		m.MatchNotes,
		m.ID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrMatchNotFound
	}
	return nil
}

func (r *matchRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM matches WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrMatchNotFound
	}
	return nil
}
