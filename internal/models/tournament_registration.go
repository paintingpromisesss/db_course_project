package models

import (
	"encoding/json"
	"time"
)

type TournamentRegistration struct {
	ID             int64           `db:"id" json:"id"`
	TournamentID   int64           `db:"tournament_id" json:"tournament_id"`
	TeamID         int64           `db:"team_id" json:"team_id"`
	SeedNumber     *int            `db:"seed_number" json:"seed_number"`
	Status         string          `db:"status" json:"status"`
	ManagerContact *string         `db:"manager_contact" json:"manager_contact"`
	RosterSnapshot json.RawMessage `db:"roster_snapshot" json:"roster_snapshot" swaggertype:"object"`
	IsInvited      bool            `db:"is_invited" json:"is_invited"`
	RegisteredAt   time.Time       `db:"registered_at" json:"registered_at"`
}

type TournamentRegistrationFilter struct {
	TournamentID *int64
	TeamID       *int64
	Status       string
	Limit        int
	Offset       int
}
