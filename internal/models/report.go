package models

import "time"

type ActiveRosterView struct {
	TeamID      int64     `db:"team_id" json:"team_id"`
	TeamName    string    `db:"team_name" json:"team_name"`
	Tag         string    `db:"tag" json:"tag"`
	PlayerID    int64     `db:"player_id" json:"player_id"`
	Nickname    string    `db:"nickname" json:"nickname"`
	CountryCode string    `db:"country_code" json:"country_code"`
	Role        string    `db:"role" json:"role"`
	JoinDate    time.Time `db:"join_date" json:"join_date"`
}

type MatchResultView struct {
	MatchID         int64     `db:"match_id" json:"match_id"`
	TournamentID    int64     `db:"tournament_id" json:"tournament_id"`
	StartTime       time.Time `db:"start_time" json:"start_time"`
	Stage           *string   `db:"stage" json:"stage"`
	Format          string    `db:"format" json:"format"`
	WinnerTeamID    *int64    `db:"winner_team_id" json:"winner_team_id"`
	GamesPlayed     int64     `db:"games_played" json:"games_played"`
	TotalScoreTeam1 *int64    `db:"total_score_team1" json:"total_score_team1"`
	TotalScoreTeam2 *int64    `db:"total_score_team2" json:"total_score_team2"`
}

type PlayerCareerStats struct {
	PlayerID int64   `db:"player_id" json:"player_id"`
	Nickname string  `db:"nickname" json:"nickname"`
	Kills    int64   `db:"kills" json:"kills"`
	Deaths   int64   `db:"deaths" json:"deaths"`
	Assists  int64   `db:"assists" json:"assists"`
	Damage   int64   `db:"damage" json:"damage"`
	Gold     int64   `db:"gold" json:"gold"`
	KDA      float64 `db:"kda" json:"kda"`
}

type TournamentStanding struct {
	TeamID        int64 `db:"team_id" json:"team_id"`
	MatchesPlayed int64 `db:"matches_played" json:"matches_played"`
	Wins          int64 `db:"wins" json:"wins"`
	Losses        int64 `db:"losses" json:"losses"`
	Forfeits      int64 `db:"forfeits" json:"forfeits"`
}
