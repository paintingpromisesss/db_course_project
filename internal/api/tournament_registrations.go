package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"db_course_project/internal/models"
	"db_course_project/internal/repository"
	"db_course_project/internal/service"
)

// TournamentRegistrationHandler manages registrations.
type TournamentRegistrationHandler struct {
	svc *service.TournamentRegistrationService
}

func NewTournamentRegistrationHandler(svc *service.TournamentRegistrationService) *TournamentRegistrationHandler {
	return &TournamentRegistrationHandler{svc: svc}
}

func (h *TournamentRegistrationHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/tournament-registrations", h.Create)
	rg.GET("/tournament-registrations", h.List)
	rg.GET("/tournament-registrations/:id", h.Get)
	rg.PUT("/tournament-registrations/:id", h.Update)
	rg.DELETE("/tournament-registrations/:id", h.Delete)
}

type tournamentRegistrationRequest struct {
	TournamentID   int64           `json:"tournament_id" binding:"required"`
	TeamID         int64           `json:"team_id" binding:"required"`
	SeedNumber     *int            `json:"seed_number"`
	Status         string          `json:"status"`
	ManagerContact *string         `json:"manager_contact"`
	RosterSnapshot json.RawMessage `json:"roster_snapshot" swaggertype:"object"`
	IsInvited      *bool           `json:"is_invited"`
}

// @Summary Create tournament registration
// @Tags TournamentRegistrations
// @Accept json
// @Produce json
// @Param payload body tournamentRegistrationRequest true "Registration payload"
// @Success 201 {object} models.TournamentRegistration
// @Failure 400 {object} map[string]string
// @Router /tournament-registrations [post]
func (h *TournamentRegistrationHandler) Create(c *gin.Context) {
	var req tournamentRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	isInvited := false
	if req.IsInvited != nil {
		isInvited = *req.IsInvited
	}
	reg := &models.TournamentRegistration{
		TournamentID:   req.TournamentID,
		TeamID:         req.TeamID,
		SeedNumber:     req.SeedNumber,
		Status:         req.Status,
		ManagerContact: req.ManagerContact,
		RosterSnapshot: req.RosterSnapshot,
		IsInvited:      isInvited,
	}
	if err := h.svc.Create(c.Request.Context(), reg); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, reg, nil)
}

// @Summary Get tournament registration
// @Tags TournamentRegistrations
// @Produce json
// @Param id path int true "Registration ID"
// @Success 200 {object} models.TournamentRegistration
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tournament-registrations/{id} [get]
func (h *TournamentRegistrationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	reg, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTournamentRegistrationNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, reg, nil)
}

// @Summary List tournament registrations
// @Tags TournamentRegistrations
// @Produce json
// @Param tournament_id query int false "Tournament ID"
// @Param team_id query int false "Team ID"
// @Param status query string false "Status"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /tournament-registrations [get]
func (h *TournamentRegistrationHandler) List(c *gin.Context) {
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
	filter := models.TournamentRegistrationFilter{
		TournamentID: tournamentID,
		TeamID:       teamID,
		Status:       c.Query("status"),
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

// @Summary Update tournament registration
// @Tags TournamentRegistrations
// @Accept json
// @Produce json
// @Param id path int true "Registration ID"
// @Param payload body tournamentRegistrationRequest true "Registration payload"
// @Success 200 {object} models.TournamentRegistration
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tournament-registrations/{id} [put]
func (h *TournamentRegistrationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req tournamentRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	isInvited := false
	if req.IsInvited != nil {
		isInvited = *req.IsInvited
	}
	reg := &models.TournamentRegistration{
		ID:             id,
		TournamentID:   req.TournamentID,
		TeamID:         req.TeamID,
		SeedNumber:     req.SeedNumber,
		Status:         req.Status,
		ManagerContact: req.ManagerContact,
		RosterSnapshot: req.RosterSnapshot,
		IsInvited:      isInvited,
	}
	if err := h.svc.Update(c.Request.Context(), reg); err != nil {
		if errors.Is(err, repository.ErrTournamentRegistrationNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, reg, nil)
}

// @Summary Delete tournament registration
// @Tags TournamentRegistrations
// @Produce json
// @Param id path int true "Registration ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tournament-registrations/{id} [delete]
func (h *TournamentRegistrationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrTournamentRegistrationNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
