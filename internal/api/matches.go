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

type MatchHandler struct {
	svc *service.MatchService
}

func NewMatchHandler(svc *service.MatchService) *MatchHandler {
	return &MatchHandler{svc: svc}
}

func (h *MatchHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/matches", h.Create)
	rg.GET("/matches", h.List)
	rg.GET("/matches/:id", h.Get)
	rg.PUT("/matches/:id", h.Update)
	rg.DELETE("/matches/:id", h.Delete)
}

type matchRequest struct {
	TournamentID int64           `json:"tournament_id" binding:"required"`
	Team1ID      *int64          `json:"team1_id"`
	Team2ID      *int64          `json:"team2_id"`
	StartTime    string          `json:"start_time" binding:"required"`
	Format       string          `json:"format"`
	Stage        *string         `json:"stage"`
	WinnerTeamID *int64          `json:"winner_team_id"`
	IsForfeit    *bool           `json:"is_forfeit"`
	MatchNotes   json.RawMessage `json:"match_notes" swaggertype:"object"`
}

func parseDateTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

// @Summary Create match
// @Tags Matches
// @Accept json
// @Produce json
// @Param payload body matchRequest true "Match payload"
// @Success 201 {object} models.Match
// @Failure 400 {object} map[string]string
// @Router /matches [post]
func (h *MatchHandler) Create(c *gin.Context) {
	var req matchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	start, err := parseDateTime(req.StartTime)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid start_time")
		return
	}
	isForfeit := false
	if req.IsForfeit != nil {
		isForfeit = *req.IsForfeit
	}
	var matchNotes *json.RawMessage
	if req.MatchNotes != nil {
		mn := req.MatchNotes
		matchNotes = &mn
	}
	m := &models.Match{
		TournamentID: req.TournamentID,
		Team1ID:      req.Team1ID,
		Team2ID:      req.Team2ID,
		StartTime:    start,
		Format:       req.Format,
		Stage:        req.Stage,
		WinnerTeamID: req.WinnerTeamID,
		IsForfeit:    isForfeit,
		MatchNotes:   matchNotes,
	}
	if err := h.svc.Create(c.Request.Context(), m); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, m, nil)
}

// @Summary Get match
// @Tags Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200 {object} models.Match
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /matches/{id} [get]
func (h *MatchHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrMatchNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, m, nil)
}

// @Summary List matches
// @Tags Matches
// @Produce json
// @Param tournament_id query int false "Tournament ID"
// @Param team_id query int false "Team ID"
// @Param stage query string false "Stage"
// @Param format query string false "Format"
// @Param from query string false "From datetime (RFC3339)"
// @Param to query string false "To datetime (RFC3339)"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /matches [get]
func (h *MatchHandler) List(c *gin.Context) {
	limit, offset := ParsePagination(c)
	var tournamentID *int64
	if v := c.Query("tournament_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			tournamentID = &parsed
		}
	}
	var teamID *int64
	if v := c.Query("team_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			teamID = &parsed
		}
	}
	var fromTime *time.Time
	if v := c.Query("from"); v != "" {
		if parsed, err := parseDateTime(v); err == nil {
			fromTime = &parsed
		}
	}
	var toTime *time.Time
	if v := c.Query("to"); v != "" {
		if parsed, err := parseDateTime(v); err == nil {
			toTime = &parsed
		}
	}
	filter := models.MatchFilter{
		TournamentID: tournamentID,
		TeamID:       teamID,
		Stage:        c.Query("stage"),
		Format:       c.Query("format"),
		From:         fromTime,
		To:           toTime,
		Limit:        limit,
		Offset:       offset,
	}
	rows, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Update match
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param payload body matchRequest true "Match payload"
// @Success 200 {object} models.Match
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /matches/{id} [put]
func (h *MatchHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req matchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	start, err := parseDateTime(req.StartTime)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid start_time")
		return
	}
	isForfeit := false
	if req.IsForfeit != nil {
		isForfeit = *req.IsForfeit
	}
	var matchNotes *json.RawMessage
	if req.MatchNotes != nil {
		mn := req.MatchNotes
		matchNotes = &mn
	}
	m := &models.Match{
		ID:           id,
		TournamentID: req.TournamentID,
		Team1ID:      req.Team1ID,
		Team2ID:      req.Team2ID,
		StartTime:    start,
		Format:       req.Format,
		Stage:        req.Stage,
		WinnerTeamID: req.WinnerTeamID,
		IsForfeit:    isForfeit,
		MatchNotes:   matchNotes,
	}
	if err := h.svc.Update(c.Request.Context(), m); err != nil {
		if errors.Is(err, repository.ErrMatchNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, m, nil)
}

// @Summary Delete match
// @Tags Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /matches/{id} [delete]
func (h *MatchHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrMatchNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
