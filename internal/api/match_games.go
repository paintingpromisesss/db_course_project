package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"db_course_project/internal/models"
	"db_course_project/internal/repository"
	"db_course_project/internal/service"
)

// MatchGameHandler manages per-map endpoints.
type MatchGameHandler struct {
	svc *service.MatchGameService
}

func NewMatchGameHandler(svc *service.MatchGameService) *MatchGameHandler {
	return &MatchGameHandler{svc: svc}
}

func (h *MatchGameHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/match-games", h.Create)
	rg.GET("/match-games", h.List)
	rg.GET("/match-games/:id", h.Get)
	rg.PUT("/match-games/:id", h.Update)
	rg.DELETE("/match-games/:id", h.Delete)
}

type matchGameRequest struct {
	MatchID           int64           `json:"match_id" binding:"required"`
	MapName           string          `json:"map_name" binding:"required"`
	GameNumber        int             `json:"game_number" binding:"required"`
	DurationSeconds   *int            `json:"duration_seconds"`
	WinnerTeamID      *int64          `json:"winner_team_id"`
	ScoreTeam1        *int            `json:"score_team1"`
	ScoreTeam2        *int            `json:"score_team2"`
	StartedAt         *string         `json:"started_at"`
	HadTechnicalPause *bool           `json:"had_technical_pause"`
	PickBanPhase      json.RawMessage `json:"pick_ban_phase" swaggertype:"object"`
}

func parseDateTimePtr(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// @Summary Create match game
// @Tags MatchGames
// @Accept json
// @Produce json
// @Param payload body matchGameRequest true "Match game payload"
// @Success 201 {object} models.MatchGame
// @Failure 400 {object} map[string]string
// @Router /match-games [post]
func (h *MatchGameHandler) Create(c *gin.Context) {
	var req matchGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	startedAt, err := parseDateTimePtr(req.StartedAt)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid started_at")
		return
	}
	hasTech := false
	if req.HadTechnicalPause != nil {
		hasTech = *req.HadTechnicalPause
	}
	g := &models.MatchGame{
		MatchID:           req.MatchID,
		MapName:           req.MapName,
		GameNumber:        req.GameNumber,
		DurationSeconds:   req.DurationSeconds,
		WinnerTeamID:      req.WinnerTeamID,
		ScoreTeam1:        req.ScoreTeam1,
		ScoreTeam2:        req.ScoreTeam2,
		StartedAt:         startedAt,
		HadTechnicalPause: hasTech,
		PickBanPhase:      req.PickBanPhase,
	}
	if err := h.svc.Create(c.Request.Context(), g); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, g, nil)
}

// @Summary Get match game
// @Tags MatchGames
// @Produce json
// @Param id path int true "Match game ID"
// @Success 200 {object} models.MatchGame
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /match-games/{id} [get]
func (h *MatchGameHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	g, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrMatchGameNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, g, nil)
}

// @Summary List match games
// @Tags MatchGames
// @Produce json
// @Param match_id query int false "Match ID"
// @Param winner_team_id query int false "Winner team ID"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /match-games [get]
func (h *MatchGameHandler) List(c *gin.Context) {
	limit, offset := ParsePagination(c)
	var matchID *int64
	if v := c.Query("match_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			matchID = &parsed
		}
	}
	var winnerID *int64
	if v := c.Query("winner_team_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			winnerID = &parsed
		}
	}
	filter := models.MatchGameFilter{MatchID: matchID, WinnerTeamID: winnerID, Limit: limit, Offset: offset}
	rows, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Update match game
// @Tags MatchGames
// @Accept json
// @Produce json
// @Param id path int true "Match game ID"
// @Param payload body matchGameRequest true "Match game payload"
// @Success 200 {object} models.MatchGame
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /match-games/{id} [put]
func (h *MatchGameHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req matchGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	startedAt, err := parseDateTimePtr(req.StartedAt)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid started_at")
		return
	}
	hasTech := false
	if req.HadTechnicalPause != nil {
		hasTech = *req.HadTechnicalPause
	}
	g := &models.MatchGame{
		ID:                id,
		MatchID:           req.MatchID,
		MapName:           req.MapName,
		GameNumber:        req.GameNumber,
		DurationSeconds:   req.DurationSeconds,
		WinnerTeamID:      req.WinnerTeamID,
		ScoreTeam1:        req.ScoreTeam1,
		ScoreTeam2:        req.ScoreTeam2,
		StartedAt:         startedAt,
		HadTechnicalPause: hasTech,
		PickBanPhase:      req.PickBanPhase,
	}
	if err := h.svc.Update(c.Request.Context(), g); err != nil {
		if errors.Is(err, repository.ErrMatchGameNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, g, nil)
}

// @Summary Delete match game
// @Tags MatchGames
// @Produce json
// @Param id path int true "Match game ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /match-games/{id} [delete]
func (h *MatchGameHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrMatchGameNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
