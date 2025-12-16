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

type TeamRepository interface {
	Create(ctx context.Context, t *models.Team) error
	GetByID(ctx context.Context, id int64) (*models.Team, error)
	List(ctx context.Context, filter models.TeamFilter) ([]models.Team, int, error)
	Update(ctx context.Context, t *models.Team) error
	Delete(ctx context.Context, id int64) error
}

func NewTeamRepository(db *sqlx.DB) TeamRepository {
	return &teamRepo{db: db}
}

var ErrTeamNotFound = errors.New("team not found")

type teamRepo struct {
	db *sqlx.DB
}

func (r *teamRepo) Create(ctx context.Context, t *models.Team) error {
	query := `INSERT INTO teams (name, tag, country_code, discipline_id, logo_url, world_ranking, is_verified)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)
			 RETURNING id, created_at`
	return r.db.QueryRowxContext(ctx, query,
		t.Name,
		t.Tag,
		t.CountryCode,
		t.DisciplineID,
		t.LogoURL,
		t.WorldRanking,
		t.IsVerified,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *teamRepo) GetByID(ctx context.Context, id int64) (*models.Team, error) {
	var t models.Team
	query := `SELECT id, name, tag, country_code, discipline_id, created_at, logo_url, world_ranking, is_verified
			  FROM teams WHERE id=$1`
	if err := r.db.GetContext(ctx, &t, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *teamRepo) List(ctx context.Context, filter models.TeamFilter) ([]models.Team, int, error) {
	base := `FROM teams WHERE 1=1`
	args := []any{}
	conds := strings.Builder{}

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := strconv.Itoa(len(args))
		conds.WriteString(` AND (LOWER(name) LIKE LOWER($` + idx + `) OR LOWER(tag) LIKE LOWER($` + idx + `))`)
	}
	if filter.CountryCode != "" {
		args = append(args, filter.CountryCode)
		conds.WriteString(` AND LOWER(country_code) = LOWER($` + strconv.Itoa(len(args)) + `)`)
	}
	if filter.DisciplineID != nil {
		args = append(args, *filter.DisciplineID)
		conds.WriteString(` AND discipline_id = $` + strconv.Itoa(len(args)))
	}
	if filter.IsVerified != nil {
		args = append(args, *filter.IsVerified)
		conds.WriteString(` AND is_verified = $` + strconv.Itoa(len(args)))
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT id, name, tag, country_code, discipline_id, created_at, logo_url, world_ranking, is_verified ` + base + conds.String() + `
				 ORDER BY name ASC, id ASC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.Team{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *teamRepo) Update(ctx context.Context, t *models.Team) error {
	query := `UPDATE teams SET name=$1, tag=$2, country_code=$3, discipline_id=$4, logo_url=$5, world_ranking=$6, is_verified=$7
			  WHERE id=$8 RETURNING created_at`
	if err := r.db.QueryRowxContext(ctx, query,
		t.Name,
		t.Tag,
		t.CountryCode,
		t.DisciplineID,
		t.LogoURL,
		t.WorldRanking,
		t.IsVerified,
		t.ID,
	).Scan(&t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTeamNotFound
		}
		return err
	}
	return nil
}

func (r *teamRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM teams WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrTeamNotFound
	}
	return nil
}
