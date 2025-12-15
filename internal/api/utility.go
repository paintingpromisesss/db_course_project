package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"db_course_project/internal/service"
)

// UtilityHandler hosts report and import endpoints.
type UtilityHandler struct {
	reports  *service.ReportService
	importer *service.ImportService
}

func NewUtilityHandler(reports *service.ReportService, importer *service.ImportService) *UtilityHandler {
	return &UtilityHandler{reports: reports, importer: importer}
}

func (h *UtilityHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/batch-import", h.BatchImport)
	rg.GET("/reports/active-rosters", h.ActiveRosters)
	rg.GET("/reports/match-results", h.MatchResults)
	rg.GET("/reports/player-career", h.PlayerCareer)
	rg.GET("/reports/tournament-standings", h.TournamentStandings)
	rg.GET("/reports/player-kda", h.PlayerKDA)
}

// @Summary Batch import players
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.PlayerImportInput true "Players to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import [post]
func (h *UtilityHandler) BatchImport(c *gin.Context) {
	var payload []service.PlayerImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportPlayers(c.Request.Context(), "players_api", payload)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, summary, nil)
}

// @Summary Active roster report
// @Tags Utility
// @Produce json
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /reports/active-rosters [get]
func (h *UtilityHandler) ActiveRosters(c *gin.Context) {
	limit, offset := ParsePagination(c)
	rows, total, err := h.reports.ActiveRosters(c.Request.Context(), limit, offset)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Match results report
// @Tags Utility
// @Produce json
// @Param tournament_id query int false "Tournament ID"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /reports/match-results [get]
func (h *UtilityHandler) MatchResults(c *gin.Context) {
	limit, offset := ParsePagination(c)
	var tournamentID *int64
	if v := c.Query("tournament_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			tournamentID = &parsed
		}
	}
	rows, total, err := h.reports.MatchResults(c.Request.Context(), tournamentID, limit, offset)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Player career report
// @Tags Utility
// @Produce json
// @Param search query string false "Search by nickname"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /reports/player-career [get]
func (h *UtilityHandler) PlayerCareer(c *gin.Context) {
	limit, offset := ParsePagination(c)
	search := c.Query("search")
	rows, total, err := h.reports.PlayerCareer(c.Request.Context(), search, limit, offset)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Tournament standings report
// @Tags Utility
// @Produce json
// @Param tournament_id query int true "Tournament ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reports/tournament-standings [get]
func (h *UtilityHandler) TournamentStandings(c *gin.Context) {
	val := c.Query("tournament_id")
	if val == "" {
		RespondError(c, http.StatusBadRequest, "tournament_id is required")
		return
	}
	tid, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid tournament_id")
		return
	}
	rows, err := h.reports.TournamentStandings(c.Request.Context(), tid)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, rows, nil)
}

// @Summary Player KDA report
// @Tags Utility
// @Produce json
// @Param player_id query int true "Player ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reports/player-kda [get]
func (h *UtilityHandler) PlayerKDA(c *gin.Context) {
	val := c.Query("player_id")
	if val == "" {
		RespondError(c, http.StatusBadRequest, "player_id is required")
		return
	}
	pid, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid player_id")
		return
	}
	kda, err := h.reports.PlayerKDA(c.Request.Context(), pid)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, gin.H{"player_id": pid, "kda": kda}, nil)
}
