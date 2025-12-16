package models

import "time"

type Player struct {
	ID          int64      `db:"id" json:"id"`
	Nickname    string     `db:"nickname" json:"nickname"`
	RealName    *string    `db:"real_name" json:"real_name"`
	CountryCode *string    `db:"country_code" json:"country_code"`
	BirthDate   *time.Time `db:"birth_date" json:"birth_date"`
	SteamID     *string    `db:"steam_id" json:"steam_id"`
	AvatarURL   *string    `db:"avatar_url" json:"avatar_url"`
	MMRRating   float64    `db:"mmr_rating" json:"mmr_rating"`
	IsRetired   bool       `db:"is_retired" json:"is_retired"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

type PlayerFilter struct {
	Search      string
	CountryCode string
	IsRetired   *bool
	MinMMR      *float64
	MaxMMR      *float64
	Limit       int
	Offset      int
}
