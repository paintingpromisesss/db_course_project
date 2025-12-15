package models

import "encoding/json"

// Discipline represents an esport discipline (game title).
type Discipline struct {
	ID          int64           `db:"id" json:"id"`
	Code        string          `db:"code" json:"code"`
	Name        string          `db:"name" json:"name"`
	Description string          `db:"description" json:"description"`
	IconURL     *string         `db:"icon_url" json:"icon_url"`
	TeamSize    *int            `db:"team_size" json:"team_size"`
	IsActive    bool            `db:"is_active" json:"is_active"`
	Metadata    json.RawMessage `db:"metadata" json:"metadata" swaggertype:"object"`
}

// DisciplineFilter holds query parameters for listing.
type DisciplineFilter struct {
	Search   string
	IsActive *bool
	Limit    int
	Offset   int
}
