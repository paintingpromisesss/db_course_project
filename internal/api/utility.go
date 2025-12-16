package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"db_course_project/internal/service"
)

type UtilityHandler struct {
	reports  *service.ReportService
	importer *service.ImportService
}

func NewUtilityHandler(reports *service.ReportService, importer *service.ImportService) *UtilityHandler {
	return &UtilityHandler{reports: reports, importer: importer}
}

func (h *UtilityHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/batch-import", h.BatchImport)
	rg.POST("/batch-import/players", h.BatchImportPlayers)
	rg.POST("/batch-import/disciplines", h.BatchImportDisciplines)
	rg.POST("/batch-import/teams", h.BatchImportTeams)
	rg.POST("/batch-import/tournaments", h.BatchImportTournaments)
	rg.POST("/batch-import/tournament-registrations", h.BatchImportTournamentRegistrations)
	rg.POST("/batch-import/matches", h.BatchImportMatches)
	rg.POST("/batch-import/match-games", h.BatchImportMatchGames)
	rg.POST("/batch-import/game-player-stats", h.BatchImportGamePlayerStats)
	rg.POST("/batch-import/squad-members", h.BatchImportSquadMembers)
	rg.POST("/batch-import/team-profiles", h.BatchImportTeamProfiles)
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

// @Summary Batch import players
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.PlayerImportInput true "Players to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/players [post]
func (h *UtilityHandler) BatchImportPlayers(c *gin.Context) {
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

// @Summary Batch import disciplines
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.DisciplineImportInput true "Disciplines to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/disciplines [post]
func (h *UtilityHandler) BatchImportDisciplines(c *gin.Context) {
	var payload []service.DisciplineImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportDisciplines(c.Request.Context(), "disciplines_api", payload)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, summary, nil)
}

// @Summary Batch import teams
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.TeamImportInput true "Teams to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/teams [post]
func (h *UtilityHandler) BatchImportTeams(c *gin.Context) {
	var payload []service.TeamImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportTeams(c.Request.Context(), "teams_api", payload)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, summary, nil)
}

// @Summary Batch import tournaments
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.TournamentImportInput true "Tournaments to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/tournaments [post]
func (h *UtilityHandler) BatchImportTournaments(c *gin.Context) {
	var payload []service.TournamentImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportTournaments(c.Request.Context(), "tournaments_api", payload)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, summary, nil)
}

// @Summary Batch import tournament registrations
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.TournamentRegistrationImportInput true "Registrations to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/tournament-registrations [post]
func (h *UtilityHandler) BatchImportTournamentRegistrations(c *gin.Context) {
	var payload []service.TournamentRegistrationImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportTournamentRegistrations(c.Request.Context(), "tournament_registrations_api", payload)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, summary, nil)
}

// @Summary Batch import matches
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.MatchImportInput true "Matches to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/matches [post]
func (h *UtilityHandler) BatchImportMatches(c *gin.Context) {
	var payload []service.MatchImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportMatches(c.Request.Context(), "matches_api", payload)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, summary, nil)
}

// @Summary Batch import match games
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.MatchGameImportInput true "Match games to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/match-games [post]
func (h *UtilityHandler) BatchImportMatchGames(c *gin.Context) {
	var payload []service.MatchGameImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportMatchGames(c.Request.Context(), "match_games_api", payload)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, summary, nil)
}

// @Summary Batch import game player stats
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.GamePlayerStatImportInput true "Game player stats to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/game-player-stats [post]
func (h *UtilityHandler) BatchImportGamePlayerStats(c *gin.Context) {
	var payload []service.GamePlayerStatImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportGamePlayerStats(c.Request.Context(), "game_player_stats_api", payload)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, summary, nil)
}

// @Summary Batch import squad members
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.SquadMemberImportInput true "Squad members to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/squad-members [post]
func (h *UtilityHandler) BatchImportSquadMembers(c *gin.Context) {
	var payload []service.SquadMemberImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportSquadMembers(c.Request.Context(), "squad_members_api", payload)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, summary, nil)
}

// @Summary Batch import team profiles
// @Tags Utility
// @Accept json
// @Produce json
// @Param payload body []service.TeamProfileImportInput true "Team profiles to import"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /batch-import/team-profiles [post]
func (h *UtilityHandler) BatchImportTeamProfiles(c *gin.Context) {
	var payload []service.TeamProfileImportInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	summary, err := h.importer.ImportTeamProfiles(c.Request.Context(), "team_profiles_api", payload)
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
