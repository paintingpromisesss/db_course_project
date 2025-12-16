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

// PlayerRepository defines persistence for players.
type PlayerRepository interface {
	Create(ctx context.Context, p *models.Player) error
	GetByID(ctx context.Context, id int64) (*models.Player, error)
	List(ctx context.Context, filter models.PlayerFilter) ([]models.Player, int, error)
	Update(ctx context.Context, p *models.Player) error
	Delete(ctx context.Context, id int64) error
}

func NewPlayerRepository(db *sqlx.DB) PlayerRepository {
	return &playerRepo{db: db}
}

// ErrPlayerNotFound signals missing row.
var ErrPlayerNotFound = errors.New("player not found")

type playerRepo struct {
	db *sqlx.DB
}

func (r *playerRepo) Create(ctx context.Context, p *models.Player) error {
	query := `INSERT INTO players (nickname, real_name, country_code, birth_date, steam_id, avatar_url, mmr_rating, is_retired)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			 RETURNING id, created_at`
	return r.db.QueryRowxContext(ctx, query,
		p.Nickname,
		p.RealName,
		p.CountryCode,
		p.BirthDate,
		p.SteamID,
		p.AvatarURL,
		p.MMRRating,
		p.IsRetired,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *playerRepo) GetByID(ctx context.Context, id int64) (*models.Player, error) {
	var p models.Player
	query := `SELECT id, nickname, real_name, country_code, birth_date, steam_id, avatar_url, mmr_rating, is_retired, created_at
			  FROM players WHERE id=$1`
	if err := r.db.GetContext(ctx, &p, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *playerRepo) List(ctx context.Context, filter models.PlayerFilter) ([]models.Player, int, error) {
	base := `FROM players WHERE 1=1`
	args := []any{}
	conds := strings.Builder{}

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := strconv.Itoa(len(args))
		conds.WriteString(` AND (LOWER(nickname) LIKE LOWER($` + idx + `) OR LOWER(real_name) LIKE LOWER($` + idx + `))`)
	}
	if filter.CountryCode != "" {
		args = append(args, filter.CountryCode)
		conds.WriteString(` AND LOWER(country_code) = LOWER($` + strconv.Itoa(len(args)) + `)`)
	}
	if filter.IsRetired != nil {
		args = append(args, *filter.IsRetired)
		conds.WriteString(` AND is_retired = $` + strconv.Itoa(len(args)))
	}
	if filter.MinMMR != nil {
		args = append(args, *filter.MinMMR)
		conds.WriteString(` AND mmr_rating >= $` + strconv.Itoa(len(args)))
	}
	if filter.MaxMMR != nil {
		args = append(args, *filter.MaxMMR)
		conds.WriteString(` AND mmr_rating <= $` + strconv.Itoa(len(args)))
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT id, nickname, real_name, country_code, birth_date, steam_id, avatar_url, mmr_rating, is_retired, created_at ` + base + conds.String() + `
				 ORDER BY nickname ASC, id ASC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.Player{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *playerRepo) Update(ctx context.Context, p *models.Player) error {
	query := `UPDATE players SET nickname=$1, real_name=$2, country_code=$3, birth_date=$4, steam_id=$5, avatar_url=$6, mmr_rating=$7, is_retired=$8
			  WHERE id=$9 RETURNING created_at`
	if err := r.db.QueryRowxContext(ctx, query,
		p.Nickname,
		p.RealName,
		p.CountryCode,
		p.BirthDate,
		p.SteamID,
		p.AvatarURL,
		p.MMRRating,
		p.IsRetired,
		p.ID,
	).Scan(&p.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPlayerNotFound
		}
		return err
	}
	return nil
}

func (r *playerRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM players WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrPlayerNotFound
	}
	return nil
}
