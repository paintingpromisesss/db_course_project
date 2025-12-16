package models

import (
	"encoding/json"
	"time"
)

type MatchGame struct {
	ID                int64           `db:"id" json:"id"`
	MatchID           int64           `db:"match_id" json:"match_id"`
	MapName           string          `db:"map_name" json:"map_name"`
	GameNumber        int             `db:"game_number" json:"game_number"`
	DurationSeconds   *int            `db:"duration_seconds" json:"duration_seconds"`
	WinnerTeamID      *int64          `db:"winner_team_id" json:"winner_team_id"`
	ScoreTeam1        *int            `db:"score_team1" json:"score_team1"`
	ScoreTeam2        *int            `db:"score_team2" json:"score_team2"`
	StartedAt         *time.Time      `db:"started_at" json:"started_at"`
	HadTechnicalPause bool            `db:"had_technical_pause" json:"had_technical_pause"`
	PickBanPhase      json.RawMessage `db:"pick_ban_phase" json:"pick_ban_phase" swaggertype:"object"`
}

type MatchGameFilter struct {
	MatchID      *int64
	WinnerTeamID *int64
	Limit        int
	Offset       int
}
