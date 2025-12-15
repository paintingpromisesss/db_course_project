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

// TournamentRegistrationRepository persistence.
type TournamentRegistrationRepository interface {
	Create(ctx context.Context, r *models.TournamentRegistration) error
	GetByID(ctx context.Context, id int64) (*models.TournamentRegistration, error)
	List(ctx context.Context, filter models.TournamentRegistrationFilter) ([]models.TournamentRegistration, int, error)
	Update(ctx context.Context, r *models.TournamentRegistration) error
	Delete(ctx context.Context, id int64) error
}

func NewTournamentRegistrationRepository(db *sqlx.DB) TournamentRegistrationRepository {
	return &tournamentRegistrationRepo{db: db}
}

// ErrTournamentRegistrationNotFound signals missing row.
var ErrTournamentRegistrationNotFound = errors.New("tournament registration not found")

type tournamentRegistrationRepo struct {
	db *sqlx.DB
}

func (r *tournamentRegistrationRepo) Create(ctx context.Context, reg *models.TournamentRegistration) error {
	query := `INSERT INTO tournament_registrations (tournament_id, team_id, seed_number, status, manager_contact, roster_snapshot, is_invited)
			  VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, registered_at`
	return r.db.QueryRowxContext(ctx, query,
		reg.TournamentID,
		reg.TeamID,
		reg.SeedNumber,
		reg.Status,
		reg.ManagerContact,
		reg.RosterSnapshot,
		reg.IsInvited,
	).Scan(&reg.ID, &reg.RegisteredAt)
}

func (r *tournamentRegistrationRepo) GetByID(ctx context.Context, id int64) (*models.TournamentRegistration, error) {
	var reg models.TournamentRegistration
	query := `SELECT id, tournament_id, team_id, seed_number, status, manager_contact, roster_snapshot, is_invited, registered_at
			  FROM tournament_registrations WHERE id=$1`
	if err := r.db.GetContext(ctx, &reg, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTournamentRegistrationNotFound
		}
		return nil, err
	}
	return &reg, nil
}

func (r *tournamentRegistrationRepo) List(ctx context.Context, filter models.TournamentRegistrationFilter) ([]models.TournamentRegistration, int, error) {
	base := `FROM tournament_registrations WHERE 1=1`
	args := []any{}
	conds := strings.Builder{}

	if filter.TournamentID != nil {
		args = append(args, *filter.TournamentID)
		conds.WriteString(` AND tournament_id = $` + strconv.Itoa(len(args)))
	}
	if filter.TeamID != nil {
		args = append(args, *filter.TeamID)
		conds.WriteString(` AND team_id = $` + strconv.Itoa(len(args)))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		conds.WriteString(` AND status = $` + strconv.Itoa(len(args)))
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT id, tournament_id, team_id, seed_number, status, manager_contact, roster_snapshot, is_invited, registered_at ` + base + conds.String() +
		` ORDER BY registered_at DESC, id DESC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.TournamentRegistration{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *tournamentRegistrationRepo) Update(ctx context.Context, reg *models.TournamentRegistration) error {
	query := `UPDATE tournament_registrations SET tournament_id=$1, team_id=$2, seed_number=$3, status=$4, manager_contact=$5, roster_snapshot=$6, is_invited=$7
			  WHERE id=$8`
	res, err := r.db.ExecContext(ctx, query,
		reg.TournamentID,
		reg.TeamID,
		reg.SeedNumber,
		reg.Status,
		reg.ManagerContact,
		reg.RosterSnapshot,
		reg.IsInvited,
		reg.ID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrTournamentRegistrationNotFound
	}
	return nil
}

func (r *tournamentRegistrationRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tournament_registrations WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrTournamentRegistrationNotFound
	}
	return nil
}
