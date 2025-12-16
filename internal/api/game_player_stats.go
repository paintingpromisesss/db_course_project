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

type GamePlayerStatHandler struct {
	svc *service.GamePlayerStatService
}

func NewGamePlayerStatHandler(svc *service.GamePlayerStatService) *GamePlayerStatHandler {
	return &GamePlayerStatHandler{svc: svc}
}

func (h *GamePlayerStatHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/game-player-stats", h.Create)
	rg.GET("/game-player-stats", h.List)
	rg.GET("/game-player-stats/:id", h.Get)
	rg.PUT("/game-player-stats/:id", h.Update)
	rg.DELETE("/game-player-stats/:id", h.Delete)
}

type gamePlayerStatRequest struct {
	GameID      int64   `json:"game_id" binding:"required"`
	PlayerID    int64   `json:"player_id" binding:"required"`
	TeamID      *int64  `json:"team_id"`
	Kills       int     `json:"kills"`
	Deaths      int     `json:"deaths"`
	Assists     int     `json:"assists"`
	HeroName    *string `json:"hero_name"`
	DamageDealt int     `json:"damage_dealt"`
	GoldEarned  int     `json:"gold_earned"`
	WasMVP      *bool   `json:"was_mvp"`
}

// @Summary Create game player stats
// @Tags GamePlayerStats
// @Accept json
// @Produce json
// @Param payload body gamePlayerStatRequest true "Game player stats payload"
// @Success 201 {object} models.GamePlayerStat
// @Failure 400 {object} map[string]string
// @Router /game-player-stats [post]
func (h *GamePlayerStatHandler) Create(c *gin.Context) {
	var req gamePlayerStatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	wasMVP := false
	if req.WasMVP != nil {
		wasMVP = *req.WasMVP
	}
	st := &models.GamePlayerStat{
		GameID:      req.GameID,
		PlayerID:    req.PlayerID,
		TeamID:      req.TeamID,
		Kills:       req.Kills,
		Deaths:      req.Deaths,
		Assists:     req.Assists,
		HeroName:    req.HeroName,
		DamageDealt: req.DamageDealt,
		GoldEarned:  req.GoldEarned,
		WasMVP:      wasMVP,
	}
	if err := h.svc.Create(c.Request.Context(), st); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, st, nil)
}

// @Summary Get game player stats
// @Tags GamePlayerStats
// @Produce json
// @Param id path int true "Stat ID"
// @Success 200 {object} models.GamePlayerStat
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /game-player-stats/{id} [get]
func (h *GamePlayerStatHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	st, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrGamePlayerStatNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, st, nil)
}

// @Summary List game player stats
// @Tags GamePlayerStats
// @Produce json
// @Param game_id query int false "Game ID"
// @Param player_id query int false "Player ID"
// @Param team_id query int false "Team ID"
// @Param was_mvp query bool false "Filter by MVP"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /game-player-stats [get]
func (h *GamePlayerStatHandler) List(c *gin.Context) {
	limit, offset := ParsePagination(c)
	var gameID *int64
	if v := c.Query("game_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			gameID = &parsed
		}
	}
	var playerID *int64
	if v := c.Query("player_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			playerID = &parsed
		}
	}
	var teamID *int64
	if v := c.Query("team_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			teamID = &parsed
		}
	}
	var wasMVP *bool
	if v := c.Query("was_mvp"); v != "" {
		switch v {
		case "true":
			b := true
			wasMVP = &b
		case "false":
			b := false
			wasMVP = &b
		}
	}
	filter := models.GamePlayerStatFilter{GameID: gameID, PlayerID: playerID, TeamID: teamID, WasMVP: wasMVP, Limit: limit, Offset: offset}
	rows, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Update game player stats
// @Tags GamePlayerStats
// @Accept json
// @Produce json
// @Param id path int true "Stat ID"
// @Param payload body gamePlayerStatRequest true "Game player stats payload"
// @Success 200 {object} models.GamePlayerStat
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /game-player-stats/{id} [put]
func (h *GamePlayerStatHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req gamePlayerStatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	wasMVP := false
	if req.WasMVP != nil {
		wasMVP = *req.WasMVP
	}
	st := &models.GamePlayerStat{
		ID:          id,
		GameID:      req.GameID,
		PlayerID:    req.PlayerID,
		TeamID:      req.TeamID,
		Kills:       req.Kills,
		Deaths:      req.Deaths,
		Assists:     req.Assists,
		HeroName:    req.HeroName,
		DamageDealt: req.DamageDealt,
		GoldEarned:  req.GoldEarned,
		WasMVP:      wasMVP,
	}
	if err := h.svc.Update(c.Request.Context(), st); err != nil {
		if errors.Is(err, repository.ErrGamePlayerStatNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, st, nil)
}

// @Summary Delete game player stats
// @Tags GamePlayerStats
// @Produce json
// @Param id path int true "Stat ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /game-player-stats/{id} [delete]
func (h *GamePlayerStatHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrGamePlayerStatNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
