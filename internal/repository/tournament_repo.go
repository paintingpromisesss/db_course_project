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

// TournamentRepository persistence.
type TournamentRepository interface {
	Create(ctx context.Context, t *models.Tournament) error
	GetByID(ctx context.Context, id int64) (*models.Tournament, error)
	List(ctx context.Context, filter models.TournamentFilter) ([]models.Tournament, int, error)
	Update(ctx context.Context, t *models.Tournament) error
	Delete(ctx context.Context, id int64) error
}

func NewTournamentRepository(db *sqlx.DB) TournamentRepository {
	return &tournamentRepo{db: db}
}

// ErrTournamentNotFound signals missing row.
var ErrTournamentNotFound = errors.New("tournament not found")

type tournamentRepo struct {
	db *sqlx.DB
}

func (r *tournamentRepo) Create(ctx context.Context, t *models.Tournament) error {
	query := `INSERT INTO tournaments (discipline_id, name, start_date, end_date, prize_pool, currency, status, is_online, bracket_config)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
			 RETURNING id`
	return r.db.QueryRowxContext(ctx, query,
		t.DisciplineID,
		t.Name,
		t.StartDate,
		t.EndDate,
		t.PrizePool,
		t.Currency,
		t.Status,
		t.IsOnline,
		t.BracketConfig,
	).Scan(&t.ID)
}

func (r *tournamentRepo) GetByID(ctx context.Context, id int64) (*models.Tournament, error) {
	var t models.Tournament
	query := `SELECT id, discipline_id, name, start_date, end_date, prize_pool, currency, status, is_online, bracket_config
			  FROM tournaments WHERE id=$1`
	if err := r.db.GetContext(ctx, &t, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTournamentNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *tournamentRepo) List(ctx context.Context, filter models.TournamentFilter) ([]models.Tournament, int, error) {
	base := `FROM tournaments WHERE 1=1`
	args := []any{}
	conds := strings.Builder{}

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := strconv.Itoa(len(args))
		conds.WriteString(` AND LOWER(name) LIKE LOWER($` + idx + `)`)
	}
	if filter.DisciplineID != nil {
		args = append(args, *filter.DisciplineID)
		conds.WriteString(` AND discipline_id = $` + strconv.Itoa(len(args)))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		conds.WriteString(` AND status = $` + strconv.Itoa(len(args)))
	}
	if filter.StartFrom != nil {
		args = append(args, *filter.StartFrom)
		conds.WriteString(` AND start_date >= $` + strconv.Itoa(len(args)))
	}
	if filter.StartTo != nil {
		args = append(args, *filter.StartTo)
		conds.WriteString(` AND start_date <= $` + strconv.Itoa(len(args)))
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT id, discipline_id, name, start_date, end_date, prize_pool, currency, status, is_online, bracket_config ` + base + conds.String() + `
				 ORDER BY start_date DESC, id DESC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.Tournament{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *tournamentRepo) Update(ctx context.Context, t *models.Tournament) error {
	query := `UPDATE tournaments SET discipline_id=$1, name=$2, start_date=$3, end_date=$4, prize_pool=$5, currency=$6, status=$7, is_online=$8, bracket_config=$9
			 WHERE id=$10`
	res, err := r.db.ExecContext(ctx, query,
		t.DisciplineID,
		t.Name,
		t.StartDate,
		t.EndDate,
		t.PrizePool,
		t.Currency,
		t.Status,
		t.IsOnline,
		t.BracketConfig,
		t.ID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrTournamentNotFound
	}
	return nil
}

func (r *tournamentRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tournaments WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrTournamentNotFound
	}
	return nil
}
