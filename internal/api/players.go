package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"db_course_project/internal/models"
	"db_course_project/internal/repository"
	"db_course_project/internal/service"
)

type PlayerHandler struct {
	svc *service.PlayerService
}

func NewPlayerHandler(svc *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

func (h *PlayerHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/players", h.Create)
	rg.GET("/players", h.List)
	rg.GET("/players/:id", h.Get)
	rg.PUT("/players/:id", h.Update)
	rg.DELETE("/players/:id", h.Delete)
}

type playerRequest struct {
	Nickname    string   `json:"nickname" binding:"required"`
	RealName    *string  `json:"real_name"`
	CountryCode *string  `json:"country_code"`
	BirthDate   *string  `json:"birth_date"`
	SteamID     *string  `json:"steam_id"`
	AvatarURL   *string  `json:"avatar_url"`
	MMRRating   *float64 `json:"mmr_rating"`
	IsRetired   *bool    `json:"is_retired"`
}

func parseDatePtr(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// @Summary Create player
// @Tags Players
// @Accept json
// @Produce json
// @Param payload body playerRequest true "Player payload"
// @Success 201 {object} PlayerResponse
// @Failure 400 {object} ErrorResponse
// @Router /players [post]
func (h *PlayerHandler) Create(c *gin.Context) {
	var req playerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	birth, err := parseDatePtr(req.BirthDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid birth_date")
		return
	}
	player := &models.Player{
		Nickname:    req.Nickname,
		RealName:    req.RealName,
		CountryCode: req.CountryCode,
		BirthDate:   birth,
		SteamID:     req.SteamID,
		AvatarURL:   req.AvatarURL,
		MMRRating:   0,
		IsRetired:   false,
	}
	if req.MMRRating != nil {
		player.MMRRating = *req.MMRRating
	}
	if req.IsRetired != nil {
		player.IsRetired = *req.IsRetired
	}
	if err := h.svc.Create(c.Request.Context(), player); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, player, nil)
}

// @Summary Get player
// @Tags Players
// @Produce json
// @Param id path int true "Player ID"
// @Success 200 {object} PlayerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /players/{id} [get]
func (h *PlayerHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	player, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrPlayerNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, player, nil)
}

// @Summary List players
// @Tags Players
// @Produce json
// @Param search query string false "Search by nickname"
// @Param country_code query string false "Country code"
// @Param is_retired query bool false "Retired flag"
// @Param min_mmr query number false "Min MMR"
// @Param max_mmr query number false "Max MMR"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} PlayerListResponse
// @Failure 500 {object} ErrorResponse
// @Router /players [get]
func (h *PlayerHandler) List(c *gin.Context) {
	limit, offset := ParsePagination(c)
	var isRetired *bool
	if v := c.Query("is_retired"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			isRetired = &parsed
		}
	}
	var minMMR *float64
	if v := c.Query("min_mmr"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			minMMR = &parsed
		}
	}
	var maxMMR *float64
	if v := c.Query("max_mmr"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			maxMMR = &parsed
		}
	}
	filter := models.PlayerFilter{
		Search:      c.Query("search"),
		CountryCode: c.Query("country_code"),
		IsRetired:   isRetired,
		MinMMR:      minMMR,
		MaxMMR:      maxMMR,
		Limit:       limit,
		Offset:      offset,
	}
	rows, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Update player
// @Tags Players
// @Accept json
// @Produce json
// @Param id path int true "Player ID"
// @Param payload body playerRequest true "Player payload"
// @Success 200 {object} PlayerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /players/{id} [put]
func (h *PlayerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req playerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	birth, err := parseDatePtr(req.BirthDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid birth_date")
		return
	}
	player := &models.Player{
		ID:          id,
		Nickname:    req.Nickname,
		RealName:    req.RealName,
		CountryCode: req.CountryCode,
		BirthDate:   birth,
		SteamID:     req.SteamID,
		AvatarURL:   req.AvatarURL,
		MMRRating:   0,
		IsRetired:   false,
	}
	if req.MMRRating != nil {
		player.MMRRating = *req.MMRRating
	}
	if req.IsRetired != nil {
		player.IsRetired = *req.IsRetired
	}
	if err := h.svc.Update(c.Request.Context(), player); err != nil {
		if errors.Is(err, repository.ErrPlayerNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, player, nil)
}

// @Summary Delete player
// @Tags Players
// @Produce json
// @Param id path int true "Player ID"
// @Success 204 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /players/{id} [delete]
func (h *PlayerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrPlayerNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
