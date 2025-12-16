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

type TournamentHandler struct {
	svc *service.TournamentService
}

func NewTournamentHandler(svc *service.TournamentService) *TournamentHandler {
	return &TournamentHandler{svc: svc}
}

func (h *TournamentHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/tournaments", h.Create)
	rg.GET("/tournaments", h.List)
	rg.GET("/tournaments/:id", h.Get)
	rg.PUT("/tournaments/:id", h.Update)
	rg.DELETE("/tournaments/:id", h.Delete)
}

type tournamentRequest struct {
	DisciplineID  int64           `json:"discipline_id" binding:"required"`
	Name          string          `json:"name" binding:"required"`
	StartDate     string          `json:"start_date" binding:"required"`
	EndDate       string          `json:"end_date" binding:"required"`
	PrizePool     float64         `json:"prize_pool"`
	Currency      string          `json:"currency"`
	Status        string          `json:"status"`
	IsOnline      *bool           `json:"is_online"`
	BracketConfig json.RawMessage `json:"bracket_config" swaggertype:"object"`
}

func parseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

// @Summary Create tournament
// @Tags Tournaments
// @Accept json
// @Produce json
// @Param payload body tournamentRequest true "Tournament payload"
// @Success 201 {object} models.Tournament
// @Failure 400 {object} map[string]string
// @Router /tournaments [post]
func (h *TournamentHandler) Create(c *gin.Context) {
	var req tournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	start, err := parseDate(req.StartDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid start_date")
		return
	}
	end, err := parseDate(req.EndDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid end_date")
		return
	}
	isOnline := false
	if req.IsOnline != nil {
		isOnline = *req.IsOnline
	}
	t := &models.Tournament{
		DisciplineID:  req.DisciplineID,
		Name:          req.Name,
		StartDate:     start,
		EndDate:       end,
		PrizePool:     req.PrizePool,
		Currency:      req.Currency,
		Status:        req.Status,
		IsOnline:      isOnline,
		BracketConfig: req.BracketConfig,
	}
	if err := h.svc.Create(c.Request.Context(), t); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, t, nil)
}

// @Summary Get tournament
// @Tags Tournaments
// @Produce json
// @Param id path int true "Tournament ID"
// @Success 200 {object} models.Tournament
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tournaments/{id} [get]
func (h *TournamentHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	t, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTournamentNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, t, nil)
}

// @Summary List tournaments
// @Tags Tournaments
// @Produce json
// @Param search query string false "Search by name"
// @Param discipline_id query int false "Discipline ID"
// @Param status query string false "Tournament status"
// @Param start_from query string false "Start date from (YYYY-MM-DD)"
// @Param start_to query string false "Start date to (YYYY-MM-DD)"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /tournaments [get]
func (h *TournamentHandler) List(c *gin.Context) {
	limit, offset := ParsePagination(c)
	var disciplineID *int64
	if v := c.Query("discipline_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			disciplineID = &parsed
		}
	}
	var startFrom *time.Time
	if v := c.Query("start_from"); v != "" {
		if parsed, err := parseDate(v); err == nil {
			startFrom = &parsed
		}
	}
	var startTo *time.Time
	if v := c.Query("start_to"); v != "" {
		if parsed, err := parseDate(v); err == nil {
			startTo = &parsed
		}
	}
	filter := models.TournamentFilter{
		Search:       c.Query("search"),
		DisciplineID: disciplineID,
		Status:       c.Query("status"),
		StartFrom:    startFrom,
		StartTo:      startTo,
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

// @Summary Update tournament
// @Tags Tournaments
// @Accept json
// @Produce json
// @Param id path int true "Tournament ID"
// @Param payload body tournamentRequest true "Tournament payload"
// @Success 200 {object} models.Tournament
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tournaments/{id} [put]
func (h *TournamentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req tournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	start, err := parseDate(req.StartDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid start_date")
		return
	}
	end, err := parseDate(req.EndDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid end_date")
		return
	}
	isOnline := false
	if req.IsOnline != nil {
		isOnline = *req.IsOnline
	}
	t := &models.Tournament{
		ID:            id,
		DisciplineID:  req.DisciplineID,
		Name:          req.Name,
		StartDate:     start,
		EndDate:       end,
		PrizePool:     req.PrizePool,
		Currency:      req.Currency,
		Status:        req.Status,
		IsOnline:      isOnline,
		BracketConfig: req.BracketConfig,
	}
	if err := h.svc.Update(c.Request.Context(), t); err != nil {
		if errors.Is(err, repository.ErrTournamentNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, t, nil)
}

// @Summary Delete tournament
// @Tags Tournaments
// @Produce json
// @Param id path int true "Tournament ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tournaments/{id} [delete]
func (h *TournamentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrTournamentNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
