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

type TeamProfileRepository interface {
	Create(ctx context.Context, p *models.TeamProfile) error
	GetByTeamID(ctx context.Context, teamID int64) (*models.TeamProfile, error)
	List(ctx context.Context, filter models.TeamProfileFilter) ([]models.TeamProfile, int, error)
	Update(ctx context.Context, p *models.TeamProfile) error
	Delete(ctx context.Context, teamID int64) error
}

func NewTeamProfileRepository(db *sqlx.DB) TeamProfileRepository {
	return &teamProfileRepo{db: db}
}

var ErrTeamProfileNotFound = errors.New("team profile not found")

type teamProfileRepo struct {
	db *sqlx.DB
}

func (r *teamProfileRepo) Create(ctx context.Context, p *models.TeamProfile) error {
	query := `INSERT INTO team_profiles (team_id, coach_name, sponsor_info, headquarters, website, contact_email)
			  VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.ExecContext(ctx, query,
		p.TeamID,
		p.CoachName,
		p.SponsorInfo,
		p.Headquarters,
		p.Website,
		p.ContactEmail,
	)
	return err
}

func (r *teamProfileRepo) GetByTeamID(ctx context.Context, teamID int64) (*models.TeamProfile, error) {
	var p models.TeamProfile
	query := `SELECT team_id, coach_name, sponsor_info, headquarters, website, contact_email
			  FROM team_profiles WHERE team_id=$1`
	if err := r.db.GetContext(ctx, &p, query, teamID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeamProfileNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *teamProfileRepo) List(ctx context.Context, filter models.TeamProfileFilter) ([]models.TeamProfile, int, error) {
	base := `FROM team_profiles WHERE 1=1`
	conds := strings.Builder{}
	args := []any{}

	if filter.TeamID != nil {
		args = append(args, *filter.TeamID)
		conds.WriteString(` AND team_id = $` + strconv.Itoa(len(args)))
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT team_id, coach_name, sponsor_info, headquarters, website, contact_email ` + base + conds.String() +
		` ORDER BY team_id DESC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.TeamProfile{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *teamProfileRepo) Update(ctx context.Context, p *models.TeamProfile) error {
	query := `UPDATE team_profiles SET coach_name=$1, sponsor_info=$2, headquarters=$3, website=$4, contact_email=$5
			  WHERE team_id=$6`
	res, err := r.db.ExecContext(ctx, query,
		p.CoachName,
		p.SponsorInfo,
		p.Headquarters,
		p.Website,
		p.ContactEmail,
		p.TeamID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrTeamProfileNotFound
	}
	return nil
}

func (r *teamProfileRepo) Delete(ctx context.Context, teamID int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM team_profiles WHERE team_id=$1`, teamID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrTeamProfileNotFound
	}
	return nil
}
