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

// SquadMemberRepository persistence.
type SquadMemberRepository interface {
	Create(ctx context.Context, m *models.SquadMember) error
	GetByID(ctx context.Context, id int64) (*models.SquadMember, error)
	List(ctx context.Context, filter models.SquadMemberFilter) ([]models.SquadMember, int, error)
	Update(ctx context.Context, m *models.SquadMember) error
	Delete(ctx context.Context, id int64) error
}

func NewSquadMemberRepository(db *sqlx.DB) SquadMemberRepository {
	return &squadMemberRepo{db: db}
}

// ErrSquadMemberNotFound signals missing row.
var ErrSquadMemberNotFound = errors.New("squad member not found")

type squadMemberRepo struct {
	db *sqlx.DB
}

func (r *squadMemberRepo) Create(ctx context.Context, m *models.SquadMember) error {
	query := `INSERT INTO squad_members (team_id, player_id, role, is_standin, join_date, contract_end_date, leave_date, salary_monthly)
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	return r.db.QueryRowxContext(ctx, query,
		m.TeamID,
		m.PlayerID,
		m.Role,
		m.IsStandin,
		m.JoinDate,
		m.ContractEndDate,
		m.LeaveDate,
		m.SalaryMonthly,
	).Scan(&m.ID)
}

func (r *squadMemberRepo) GetByID(ctx context.Context, id int64) (*models.SquadMember, error) {
	var m models.SquadMember
	query := `SELECT id, team_id, player_id, role, is_standin, join_date, contract_end_date, leave_date, salary_monthly
			  FROM squad_members WHERE id=$1`
	if err := r.db.GetContext(ctx, &m, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSquadMemberNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *squadMemberRepo) List(ctx context.Context, filter models.SquadMemberFilter) ([]models.SquadMember, int, error) {
	base := `FROM squad_members WHERE 1=1`
	conds := strings.Builder{}
	args := []any{}

	if filter.TeamID != nil {
		args = append(args, *filter.TeamID)
		conds.WriteString(` AND team_id = $` + strconv.Itoa(len(args)))
	}
	if filter.PlayerID != nil {
		args = append(args, *filter.PlayerID)
		conds.WriteString(` AND player_id = $` + strconv.Itoa(len(args)))
	}
	if filter.ActiveOnly {
		conds.WriteString(` AND leave_date IS NULL`)
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT id, team_id, player_id, role, is_standin, join_date, contract_end_date, leave_date, salary_monthly ` + base + conds.String() +
		` ORDER BY join_date DESC, id DESC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.SquadMember{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *squadMemberRepo) Update(ctx context.Context, m *models.SquadMember) error {
	query := `UPDATE squad_members SET team_id=$1, player_id=$2, role=$3, is_standin=$4, join_date=$5, contract_end_date=$6, leave_date=$7, salary_monthly=$8
			  WHERE id=$9`
	res, err := r.db.ExecContext(ctx, query,
		m.TeamID,
		m.PlayerID,
		m.Role,
		m.IsStandin,
		m.JoinDate,
		m.ContractEndDate,
		m.LeaveDate,
		m.SalaryMonthly,
		m.ID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrSquadMemberNotFound
	}
	return nil
}

func (r *squadMemberRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM squad_members WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrSquadMemberNotFound
	}
	return nil
}
