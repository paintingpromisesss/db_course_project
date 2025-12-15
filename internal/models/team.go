package models

import "time"

// Team represents an esports team.
type Team struct {
	ID           int64     `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Tag          string    `db:"tag" json:"tag"`
	CountryCode  string    `db:"country_code" json:"country_code"`
	DisciplineID int64     `db:"discipline_id" json:"discipline_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	LogoURL      *string   `db:"logo_url" json:"logo_url"`
	WorldRanking float64   `db:"world_ranking" json:"world_ranking"`
	IsVerified   bool      `db:"is_verified" json:"is_verified"`
}

// TeamFilter describes list filters for teams.
type TeamFilter struct {
	Search       string
	CountryCode  string
	DisciplineID *int64
	IsVerified   *bool
	Limit        int
	Offset       int
}
