package models

// GamePlayerStat holds per-player stats for a map.
type GamePlayerStat struct {
	ID          int64   `db:"id" json:"id"`
	GameID      int64   `db:"game_id" json:"game_id"`
	PlayerID    int64   `db:"player_id" json:"player_id"`
	TeamID      *int64  `db:"team_id" json:"team_id"`
	Kills       int     `db:"kills" json:"kills"`
	Deaths      int     `db:"deaths" json:"deaths"`
	Assists     int     `db:"assists" json:"assists"`
	HeroName    *string `db:"hero_name" json:"hero_name"`
	DamageDealt int     `db:"damage_dealt" json:"damage_dealt"`
	GoldEarned  int     `db:"gold_earned" json:"gold_earned"`
	KDARatio    float64 `db:"kda_ratio" json:"kda_ratio"`
	WasMVP      bool    `db:"was_mvp" json:"was_mvp"`
}

// GamePlayerStatFilter supports listing.
type GamePlayerStatFilter struct {
	GameID   *int64
	PlayerID *int64
	TeamID   *int64
	WasMVP   *bool
	Limit    int
	Offset   int
}
