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

type TeamHandler struct {
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/teams", h.Create)
	rg.GET("/teams", h.List)
	rg.GET("/teams/:id", h.Get)
	rg.PUT("/teams/:id", h.Update)
	rg.DELETE("/teams/:id", h.Delete)
}

type teamRequest struct {
	Name         string   `json:"name" binding:"required"`
	Tag          string   `json:"tag" binding:"required"`
	CountryCode  string   `json:"country_code" binding:"required"`
	DisciplineID int64    `json:"discipline_id" binding:"required"`
	LogoURL      *string  `json:"logo_url"`
	WorldRanking *float64 `json:"world_ranking"`
	IsVerified   *bool    `json:"is_verified"`
}

// @Summary Create team
// @Tags Teams
// @Accept json
// @Produce json
// @Param payload body teamRequest true "Team payload"
// @Success 201 {object} TeamResponse
// @Failure 400 {object} ErrorResponse
// @Router /teams [post]
func (h *TeamHandler) Create(c *gin.Context) {
	var req teamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	team := &models.Team{
		Name:         req.Name,
		Tag:          req.Tag,
		CountryCode:  req.CountryCode,
		DisciplineID: req.DisciplineID,
		LogoURL:      req.LogoURL,
		WorldRanking: 0,
		IsVerified:   false,
	}
	if req.WorldRanking != nil {
		team.WorldRanking = *req.WorldRanking
	}
	if req.IsVerified != nil {
		team.IsVerified = *req.IsVerified
	}
	if err := h.svc.Create(c.Request.Context(), team); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, team, nil)
}

// @Summary Get team
// @Tags Teams
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {object} TeamResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /teams/{id} [get]
func (h *TeamHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	team, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, team, nil)
}

// @Summary List teams
// @Tags Teams
// @Produce json
// @Param search query string false "Search by name or tag"
// @Param country_code query string false "Country code"
// @Param discipline_id query int false "Discipline ID"
// @Param is_verified query bool false "Verification flag"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} TeamListResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams [get]
func (h *TeamHandler) List(c *gin.Context) {
	limit, offset := ParsePagination(c)
	var disciplineID *int64
	if v := c.Query("discipline_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			disciplineID = &parsed
		}
	}
	var isVerified *bool
	if v := c.Query("is_verified"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			isVerified = &parsed
		}
	}
	filter := models.TeamFilter{
		Search:       c.Query("search"),
		CountryCode:  c.Query("country_code"),
		DisciplineID: disciplineID,
		IsVerified:   isVerified,
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

// @Summary Update team
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Param payload body teamRequest true "Team payload"
// @Success 200 {object} TeamResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /teams/{id} [put]
func (h *TeamHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req teamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	team := &models.Team{
		ID:           id,
		Name:         req.Name,
		Tag:          req.Tag,
		CountryCode:  req.CountryCode,
		DisciplineID: req.DisciplineID,
		LogoURL:      req.LogoURL,
		WorldRanking: 0,
		IsVerified:   false,
	}
	if req.WorldRanking != nil {
		team.WorldRanking = *req.WorldRanking
	}
	if req.IsVerified != nil {
		team.IsVerified = *req.IsVerified
	}
	if err := h.svc.Update(c.Request.Context(), team); err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, team, nil)
}

// @Summary Delete team
// @Tags Teams
// @Produce json
// @Param id path int true "Team ID"
// @Success 204 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /teams/{id} [delete]
func (h *TeamHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
