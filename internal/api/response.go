package api

import (
	"github.com/gin-gonic/gin"

	"db_course_project/internal/models"
	"db_course_project/internal/service"
)

type PaginationMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// swagger:model
type ErrorDetail struct {
	Message string `json:"message"`
}

// swagger:model
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// swagger:model
type EmptyResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

// swagger:model
type DisciplineResponse struct {
	Data models.Discipline `json:"data"`
	Meta interface{}       `json:"meta"`
}

// swagger:model
type DisciplineListResponse struct {
	Data []models.Discipline `json:"data"`
	Meta PaginationMeta      `json:"meta"`
}

// swagger:model
type PlayerResponse struct {
	Data models.Player `json:"data"`
	Meta interface{}   `json:"meta"`
}

// swagger:model
type PlayerListResponse struct {
	Data []models.Player `json:"data"`
	Meta PaginationMeta  `json:"meta"`
}

// swagger:model
type TeamResponse struct {
	Data models.Team `json:"data"`
	Meta interface{} `json:"meta"`
}

// swagger:model
type TeamListResponse struct {
	Data []models.Team  `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// swagger:model
type TournamentResponse struct {
	Data models.Tournament `json:"data"`
	Meta interface{}       `json:"meta"`
}

// swagger:model
type TournamentListResponse struct {
	Data []models.Tournament `json:"data"`
	Meta PaginationMeta      `json:"meta"`
}

// swagger:model
type TournamentRegistrationResponse struct {
	Data models.TournamentRegistration `json:"data"`
	Meta interface{}                   `json:"meta"`
}

// swagger:model
type TournamentRegistrationListResponse struct {
	Data []models.TournamentRegistration `json:"data"`
	Meta PaginationMeta                  `json:"meta"`
}

// swagger:model
type MatchResponse struct {
	Data models.Match `json:"data"`
	Meta interface{}  `json:"meta"`
}

// swagger:model
type MatchListResponse struct {
	Data []models.Match `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// swagger:model
type MatchGameResponse struct {
	Data models.MatchGame `json:"data"`
	Meta interface{}      `json:"meta"`
}

// swagger:model
type MatchGameListResponse struct {
	Data []models.MatchGame `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

// swagger:model
type GamePlayerStatResponse struct {
	Data models.GamePlayerStat `json:"data"`
	Meta interface{}           `json:"meta"`
}

// swagger:model
type GamePlayerStatListResponse struct {
	Data []models.GamePlayerStat `json:"data"`
	Meta PaginationMeta          `json:"meta"`
}

// swagger:model
type SquadMemberResponse struct {
	Data models.SquadMember `json:"data"`
	Meta interface{}        `json:"meta"`
}

// swagger:model
type SquadMemberListResponse struct {
	Data []models.SquadMember `json:"data"`
	Meta PaginationMeta       `json:"meta"`
}

// swagger:model
type TeamProfileResponse struct {
	Data models.TeamProfile `json:"data"`
	Meta interface{}        `json:"meta"`
}

// swagger:model
type TeamProfileListResponse struct {
	Data []models.TeamProfile `json:"data"`
	Meta PaginationMeta       `json:"meta"`
}

// swagger:model
type ImportSummaryResponse struct {
	Data service.ImportSummary `json:"data"`
	Meta interface{}           `json:"meta"`
}

// swagger:model
type ActiveRostersResponse struct {
	Data []models.ActiveRosterView `json:"data"`
	Meta PaginationMeta            `json:"meta"`
}

// swagger:model
type MatchResultsResponse struct {
	Data []models.MatchResultView `json:"data"`
	Meta PaginationMeta           `json:"meta"`
}

// swagger:model
type PlayerCareerResponse struct {
	Data []models.PlayerCareerStats `json:"data"`
	Meta PaginationMeta             `json:"meta"`
}

// swagger:model
type TournamentStandingsResponse struct {
	Data []models.TournamentStanding `json:"data"`
	Meta interface{}                 `json:"meta"`
}

// swagger:model
type PlayerKDAData struct {
	PlayerID int64   `json:"player_id"`
	KDA      float64 `json:"kda"`
}

// swagger:model
type PlayerKDAResponse struct {
	Data PlayerKDAData `json:"data"`
	Meta interface{}   `json:"meta"`
}

func RespondData(c *gin.Context, status int, data any, meta any) {
	c.JSON(status, gin.H{
		"data": data,
		"meta": meta,
	})
}

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": message,
		},
	})
}
