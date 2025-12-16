package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"

	"db_course_project/internal/models"
)

type DisciplineRepository interface {
	Create(ctx context.Context, d *models.Discipline) error
	GetByID(ctx context.Context, id int64) (*models.Discipline, error)
	List(ctx context.Context, filter models.DisciplineFilter) ([]models.Discipline, int, error)
	Update(ctx context.Context, d *models.Discipline) error
	Delete(ctx context.Context, id int64) error
}

func NewDisciplineRepository(db *sqlx.DB) DisciplineRepository {
	return &disciplineRepo{db: db}
}

var ErrDisciplineNotFound = errors.New("discipline not found")

type disciplineRepo struct {
	db *sqlx.DB
}

func (r *disciplineRepo) Create(ctx context.Context, d *models.Discipline) error {
	query := `INSERT INTO disciplines (code, name, description, icon_url, team_size, is_active, metadata)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowxContext(ctx, query, d.Code, d.Name, d.Description, d.IconURL, d.TeamSize, d.IsActive, d.Metadata).
		Scan(&d.ID)
}

func (r *disciplineRepo) GetByID(ctx context.Context, id int64) (*models.Discipline, error) {
	var d models.Discipline
	query := `SELECT id, code, name, description, icon_url, team_size, is_active, metadata
			  FROM disciplines WHERE id = $1`
	if err := r.db.GetContext(ctx, &d, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDisciplineNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *disciplineRepo) List(ctx context.Context, filter models.DisciplineFilter) ([]models.Discipline, int, error) {
	base := `FROM disciplines WHERE 1=1`
	args := []any{}
	conditions := ""

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := strconv.Itoa(len(args))
		conditions += ` AND (LOWER(code) LIKE LOWER($` + idx + `) OR LOWER(name) LIKE LOWER($` + idx + `))`
	}
	if filter.IsActive != nil {
		args = append(args, *filter.IsActive)
		conditions += ` AND is_active = $` + strconv.Itoa(len(args))
	}

	countQuery := `SELECT count(*) ` + base + conditions
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT id, code, name, description, icon_url, team_size, is_active, metadata ` + base + conditions + `
				  ORDER BY name ASC, id ASC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	var rows []models.Discipline
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *disciplineRepo) Update(ctx context.Context, d *models.Discipline) error {
	query := `UPDATE disciplines
			  SET code=$1, name=$2, description=$3, icon_url=$4, team_size=$5, is_active=$6, metadata=$7
			  WHERE id=$8`
	res, err := r.db.ExecContext(ctx, query, d.Code, d.Name, d.Description, d.IconURL, d.TeamSize, d.IsActive, d.Metadata, d.ID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrDisciplineNotFound
	}
	return nil
}

func (r *disciplineRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM disciplines WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrDisciplineNotFound
	}
	return nil
}
