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

// SquadMemberHandler manages roster endpoints.
type SquadMemberHandler struct {
	svc *service.SquadMemberService
}

func NewSquadMemberHandler(svc *service.SquadMemberService) *SquadMemberHandler {
	return &SquadMemberHandler{svc: svc}
}

func (h *SquadMemberHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/squad-members", h.Create)
	rg.GET("/squad-members", h.List)
	rg.GET("/squad-members/:id", h.Get)
	rg.PUT("/squad-members/:id", h.Update)
	rg.DELETE("/squad-members/:id", h.Delete)
}

type squadMemberRequest struct {
	TeamID          int64    `json:"team_id" binding:"required"`
	PlayerID        int64    `json:"player_id" binding:"required"`
	Role            string   `json:"role"`
	IsStandin       *bool    `json:"is_standin"`
	JoinDate        string   `json:"join_date"`
	ContractEndDate *string  `json:"contract_end_date"`
	LeaveDate       *string  `json:"leave_date"`
	SalaryMonthly   *float64 `json:"salary_monthly"`
}

// parseDatePtrSM parses optional YYYY-MM-DD date strings.
func parseDatePtrSM(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseDateValue(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

// @Summary Add player to squad
// @Tags SquadMembers
// @Accept json
// @Produce json
// @Param payload body squadMemberRequest true "Squad member payload"
// @Success 201 {object} models.SquadMember
// @Failure 400 {object} map[string]string
// @Router /squad-members [post]
func (h *SquadMemberHandler) Create(c *gin.Context) {
	var req squadMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	joinDate := time.Now().UTC().Truncate(24 * time.Hour)
	if req.JoinDate != "" {
		if parsed, err := parseDateValue(req.JoinDate); err == nil {
			joinDate = parsed
		} else {
			RespondError(c, http.StatusBadRequest, "invalid join_date")
			return
		}
	}
	contractEnd, err := parseDatePtrSM(req.ContractEndDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid contract_end_date")
		return
	}
	leaveDate, err := parseDatePtrSM(req.LeaveDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid leave_date")
		return
	}
	isStandin := false
	if req.IsStandin != nil {
		isStandin = *req.IsStandin
	}
	m := &models.SquadMember{
		TeamID:          req.TeamID,
		PlayerID:        req.PlayerID,
		Role:            req.Role,
		IsStandin:       isStandin,
		JoinDate:        joinDate,
		ContractEndDate: contractEnd,
		LeaveDate:       leaveDate,
		SalaryMonthly:   req.SalaryMonthly,
	}
	if err := h.svc.Create(c.Request.Context(), m); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, m, nil)
}

// @Summary Get squad member
// @Tags SquadMembers
// @Produce json
// @Param id path int true "Squad member ID"
// @Success 200 {object} models.SquadMember
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /squad-members/{id} [get]
func (h *SquadMemberHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrSquadMemberNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, m, nil)
}

// @Summary List squad members
// @Tags SquadMembers
// @Produce json
// @Param team_id query int false "Team ID"
// @Param player_id query int false "Player ID"
// @Param active_only query bool false "Only active"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /squad-members [get]
func (h *SquadMemberHandler) List(c *gin.Context) {
	limit, offset := ParsePagination(c)
	var teamID *int64
	if v := c.Query("team_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			teamID = &parsed
		}
	}
	var playerID *int64
	if v := c.Query("player_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			playerID = &parsed
		}
	}
	activeOnly := c.Query("active_only") == "true"
	filter := models.SquadMemberFilter{TeamID: teamID, PlayerID: playerID, ActiveOnly: activeOnly, Limit: limit, Offset: offset}
	rows, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Update squad member
// @Tags SquadMembers
// @Accept json
// @Produce json
// @Param id path int true "Squad member ID"
// @Param payload body squadMemberRequest true "Squad member payload"
// @Success 200 {object} models.SquadMember
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /squad-members/{id} [put]
func (h *SquadMemberHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req squadMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.JoinDate == "" {
		RespondError(c, http.StatusBadRequest, "join_date is required")
		return
	}
	joinDate, err := parseDateValue(req.JoinDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid join_date")
		return
	}
	contractEnd, err := parseDatePtrSM(req.ContractEndDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid contract_end_date")
		return
	}
	leaveDate, err := parseDatePtrSM(req.LeaveDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid leave_date")
		return
	}
	isStandin := false
	if req.IsStandin != nil {
		isStandin = *req.IsStandin
	}
	m := &models.SquadMember{
		ID:              id,
		TeamID:          req.TeamID,
		PlayerID:        req.PlayerID,
		Role:            req.Role,
		IsStandin:       isStandin,
		JoinDate:        joinDate,
		ContractEndDate: contractEnd,
		LeaveDate:       leaveDate,
		SalaryMonthly:   req.SalaryMonthly,
	}
	if err := h.svc.Update(c.Request.Context(), m); err != nil {
		if errors.Is(err, repository.ErrSquadMemberNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, m, nil)
}

// @Summary Remove squad member
// @Tags SquadMembers
// @Produce json
// @Param id path int true "Squad member ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /squad-members/{id} [delete]
func (h *SquadMemberHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrSquadMemberNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
