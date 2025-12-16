package models

import (
	"encoding/json"
	"time"
)

// Match represents a series between teams.
type Match struct {
	ID           int64            `db:"id" json:"id"`
	TournamentID int64            `db:"tournament_id" json:"tournament_id"`
	Team1ID      *int64           `db:"team1_id" json:"team1_id"`
	Team2ID      *int64           `db:"team2_id" json:"team2_id"`
	StartTime    time.Time        `db:"start_time" json:"start_time"`
	Format       string           `db:"format" json:"format"`
	Stage        *string          `db:"stage" json:"stage"`
	WinnerTeamID *int64           `db:"winner_team_id" json:"winner_team_id"`
	IsForfeit    bool             `db:"is_forfeit" json:"is_forfeit"`
	MatchNotes   *json.RawMessage `db:"match_notes" json:"match_notes" swaggertype:"object"`
}

// MatchFilter supports listing.
type MatchFilter struct {
	TournamentID *int64
	TeamID       *int64
	Stage        string
	Format       string
	From         *time.Time
	To           *time.Time
	Limit        int
	Offset       int
}
