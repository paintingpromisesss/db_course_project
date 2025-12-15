package models

import (
	"encoding/json"
	"time"
)

// Tournament represents a tournament event.
type Tournament struct {
	ID            int64           `db:"id" json:"id"`
	DisciplineID  int64           `db:"discipline_id" json:"discipline_id"`
	Name          string          `db:"name" json:"name"`
	StartDate     time.Time       `db:"start_date" json:"start_date"`
	EndDate       time.Time       `db:"end_date" json:"end_date"`
	PrizePool     float64         `db:"prize_pool" json:"prize_pool"`
	Currency      string          `db:"currency" json:"currency"`
	Status        string          `db:"status" json:"status"`
	IsOnline      bool            `db:"is_online" json:"is_online"`
	BracketConfig json.RawMessage `db:"bracket_config" json:"bracket_config" swaggertype:"object"`
}

// TournamentFilter supports listing.
type TournamentFilter struct {
	Search       string
	DisciplineID *int64
	Status       string
	StartFrom    *time.Time
	StartTo      *time.Time
	Limit        int
	Offset       int
}
