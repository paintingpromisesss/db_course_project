package models

type TeamProfile struct {
	TeamID       int64   `db:"team_id" json:"team_id"`
	CoachName    *string `db:"coach_name" json:"coach_name"`
	SponsorInfo  *string `db:"sponsor_info" json:"sponsor_info"`
	Headquarters *string `db:"headquarters" json:"headquarters"`
	Website      *string `db:"website" json:"website"`
	ContactEmail *string `db:"contact_email" json:"contact_email"`
}

type TeamProfileFilter struct {
	TeamID *int64
	Limit  int
	Offset int
}
