package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"db_course_project/internal/models"
	"db_course_project/internal/repository"
	"db_course_project/internal/service"
)

// TeamProfileHandler manages team profile endpoints.
type TeamProfileHandler struct {
	svc *service.TeamProfileService
}

func NewTeamProfileHandler(svc *service.TeamProfileService) *TeamProfileHandler {
	return &TeamProfileHandler{svc: svc}
}

func (h *TeamProfileHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/team-profiles", h.Create)
	rg.GET("/team-profiles", h.List)
	rg.GET("/team-profiles/:team_id", h.Get)
	rg.PUT("/team-profiles/:team_id", h.Update)
	rg.DELETE("/team-profiles/:team_id", h.Delete)
}

type teamProfileRequest struct {
	TeamID       int64   `json:"team_id" binding:"required"`
	CoachName    *string `json:"coach_name"`
	SponsorInfo  *string `json:"sponsor_info"`
	Headquarters *string `json:"headquarters"`
	Website      *string `json:"website"`
	ContactEmail *string `json:"contact_email"`
}

// @Summary Create team profile
// @Tags TeamProfiles
// @Accept json
// @Produce json
// @Param payload body teamProfileRequest true "Team profile payload"
// @Success 201 {object} models.TeamProfile
// @Failure 400 {object} map[string]string
// @Router /team-profiles [post]
func (h *TeamProfileHandler) Create(c *gin.Context) {
	var req teamProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	p := &models.TeamProfile{
		TeamID:       req.TeamID,
		CoachName:    req.CoachName,
		SponsorInfo:  req.SponsorInfo,
		Headquarters: req.Headquarters,
		Website:      req.Website,
		ContactEmail: req.ContactEmail,
	}
	if err := h.svc.Create(c.Request.Context(), p); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, p, nil)
}

// @Summary Get team profile
// @Tags TeamProfiles
// @Produce json
// @Param team_id path int true "Team ID"
// @Success 200 {object} models.TeamProfile
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /team-profiles/{team_id} [get]
func (h *TeamProfileHandler) Get(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("team_id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid team_id")
		return
	}
	p, err := h.svc.Get(c.Request.Context(), teamID)
	if err != nil {
		if errors.Is(err, repository.ErrTeamProfileNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, p, nil)
}

// @Summary List team profiles
// @Tags TeamProfiles
// @Produce json
// @Param team_id query int false "Team ID"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /team-profiles [get]
func (h *TeamProfileHandler) List(c *gin.Context) {
	limit, offset := ParsePagination(c)
	var teamID *int64
	if v := c.Query("team_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			teamID = &parsed
		}
	}
	filter := models.TeamProfileFilter{TeamID: teamID, Limit: limit, Offset: offset}
	rows, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Update team profile
// @Tags TeamProfiles
// @Accept json
// @Produce json
// @Param team_id path int true "Team ID"
// @Param payload body teamProfileRequest true "Team profile payload"
// @Success 200 {object} models.TeamProfile
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /team-profiles/{team_id} [put]
func (h *TeamProfileHandler) Update(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("team_id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid team_id")
		return
	}
	var req teamProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	p := &models.TeamProfile{
		TeamID:       teamID,
		CoachName:    req.CoachName,
		SponsorInfo:  req.SponsorInfo,
		Headquarters: req.Headquarters,
		Website:      req.Website,
		ContactEmail: req.ContactEmail,
	}
	if err := h.svc.Update(c.Request.Context(), p); err != nil {
		if errors.Is(err, repository.ErrTeamProfileNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, p, nil)
}

// @Summary Delete team profile
// @Tags TeamProfiles
// @Produce json
// @Param team_id path int true "Team ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /team-profiles/{team_id} [delete]
func (h *TeamProfileHandler) Delete(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("team_id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid team_id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), teamID); err != nil {
		if errors.Is(err, repository.ErrTeamProfileNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
