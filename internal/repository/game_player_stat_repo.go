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

// GamePlayerStatRepository persistence.
type GamePlayerStatRepository interface {
	Create(ctx context.Context, s *models.GamePlayerStat) error
	GetByID(ctx context.Context, id int64) (*models.GamePlayerStat, error)
	List(ctx context.Context, filter models.GamePlayerStatFilter) ([]models.GamePlayerStat, int, error)
	Update(ctx context.Context, s *models.GamePlayerStat) error
	Delete(ctx context.Context, id int64) error
}

func NewGamePlayerStatRepository(db *sqlx.DB) GamePlayerStatRepository {
	return &gamePlayerStatRepo{db: db}
}

// ErrGamePlayerStatNotFound signals missing row.
var ErrGamePlayerStatNotFound = errors.New("game player stat not found")

type gamePlayerStatRepo struct {
	db *sqlx.DB
}

func (r *gamePlayerStatRepo) Create(ctx context.Context, s *models.GamePlayerStat) error {
	query := `INSERT INTO game_player_stats (game_id, player_id, team_id, kills, deaths, assists, hero_name, damage_dealt, gold_earned, was_mvp)
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, kda_ratio`
	return r.db.QueryRowxContext(ctx, query,
		s.GameID,
		s.PlayerID,
		s.TeamID,
		s.Kills,
		s.Deaths,
		s.Assists,
		s.HeroName,
		s.DamageDealt,
		s.GoldEarned,
		s.WasMVP,
	).Scan(&s.ID, &s.KDARatio)
}

func (r *gamePlayerStatRepo) GetByID(ctx context.Context, id int64) (*models.GamePlayerStat, error) {
	var s models.GamePlayerStat
	query := `SELECT id, game_id, player_id, team_id, kills, deaths, assists, hero_name, damage_dealt, gold_earned, kda_ratio, was_mvp
			  FROM game_player_stats WHERE id=$1`
	if err := r.db.GetContext(ctx, &s, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGamePlayerStatNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *gamePlayerStatRepo) List(ctx context.Context, filter models.GamePlayerStatFilter) ([]models.GamePlayerStat, int, error) {
	base := `FROM game_player_stats WHERE 1=1`
	args := []any{}
	conds := strings.Builder{}

	if filter.GameID != nil {
		args = append(args, *filter.GameID)
		conds.WriteString(` AND game_id = $` + strconv.Itoa(len(args)))
	}
	if filter.PlayerID != nil {
		args = append(args, *filter.PlayerID)
		conds.WriteString(` AND player_id = $` + strconv.Itoa(len(args)))
	}
	if filter.TeamID != nil {
		args = append(args, *filter.TeamID)
		conds.WriteString(` AND team_id = $` + strconv.Itoa(len(args)))
	}
	if filter.WasMVP != nil {
		args = append(args, *filter.WasMVP)
		conds.WriteString(` AND was_mvp = $` + strconv.Itoa(len(args)))
	}

	countQuery := `SELECT count(*) ` + base + conds.String()
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := strconv.Itoa(len(args) - 1)
	offsetIdx := strconv.Itoa(len(args))
	listQuery := `SELECT id, game_id, player_id, team_id, kills, deaths, assists, hero_name, damage_dealt, gold_earned, kda_ratio, was_mvp ` + base + conds.String() +
		` ORDER BY game_id DESC, id DESC LIMIT $` + limitIdx + ` OFFSET $` + offsetIdx

	rows := []models.GamePlayerStat{}
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *gamePlayerStatRepo) Update(ctx context.Context, s *models.GamePlayerStat) error {
	query := `UPDATE game_player_stats SET game_id=$1, player_id=$2, team_id=$3, kills=$4, deaths=$5, assists=$6, hero_name=$7, damage_dealt=$8, gold_earned=$9, was_mvp=$10
			  WHERE id=$11 RETURNING kda_ratio`
	if err := r.db.QueryRowxContext(ctx, query,
		s.GameID,
		s.PlayerID,
		s.TeamID,
		s.Kills,
		s.Deaths,
		s.Assists,
		s.HeroName,
		s.DamageDealt,
		s.GoldEarned,
		s.WasMVP,
		s.ID,
	).Scan(&s.KDARatio); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGamePlayerStatNotFound
		}
		return err
	}
	return nil
}

func (r *gamePlayerStatRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM game_player_stats WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrGamePlayerStatNotFound
	}
	return nil
}
