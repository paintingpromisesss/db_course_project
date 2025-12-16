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

type DisciplineHandler struct {
	svc *service.DisciplineService
}

func NewDisciplineHandler(svc *service.DisciplineService) *DisciplineHandler {
	return &DisciplineHandler{svc: svc}
}

type disciplineRequest struct {
	Code        string          `json:"code" binding:"required"`
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	IconURL     *string         `json:"icon_url"`
	TeamSize    *int            `json:"team_size"`
	Metadata    json.RawMessage `json:"metadata" swaggertype:"object"`
	IsActive    *bool           `json:"is_active"`
}

func (h *DisciplineHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/disciplines", h.Create)
	rg.GET("/disciplines", h.List)
	rg.GET("/disciplines/:id", h.Get)
	rg.PUT("/disciplines/:id", h.Update)
	rg.DELETE("/disciplines/:id", h.Delete)
}

// @Summary Create discipline
// @Tags Disciplines
// @Accept json
// @Produce json
// @Param payload body disciplineRequest true "Discipline payload"
// @Success 201 {object} models.Discipline
// @Failure 400 {object} map[string]string
// @Router /disciplines [post]
func (h *DisciplineHandler) Create(c *gin.Context) {
	var req disciplineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	d := &models.Discipline{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		TeamSize:    req.TeamSize,
		Metadata:    req.Metadata,
	}
	if req.IsActive != nil {
		d.IsActive = *req.IsActive
	} else {
		d.IsActive = true
	}
	if len(d.Metadata) == 0 {
		d.Metadata = json.RawMessage(`{}`)
	}

	if err := h.svc.Create(c.Request.Context(), d); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusCreated, d, nil)
}

// @Summary Get discipline
// @Tags Disciplines
// @Produce json
// @Param id path int true "Discipline ID"
// @Success 200 {object} models.Discipline
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /disciplines/{id} [get]
func (h *DisciplineHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	d, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrDisciplineNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusOK, d, nil)
}

// @Summary List disciplines
// @Tags Disciplines
// @Produce json
// @Param search query string false "Search by name or code"
// @Param is_active query bool false "Filter by active flag"
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /disciplines [get]
func (h *DisciplineHandler) List(c *gin.Context) {
	limit, offset := ParsePagination(c)

	var isActive *bool
	if v := c.Query("is_active"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err == nil {
			isActive = &parsed
		}
	}

	filter := models.DisciplineFilter{
		Search:   c.Query("search"),
		IsActive: isActive,
		Limit:    limit,
		Offset:   offset,
	}

	rows, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := PaginationMeta{Total: total, Limit: limit, Offset: offset}
	RespondData(c, http.StatusOK, rows, meta)
}

// @Summary Update discipline
// @Tags Disciplines
// @Accept json
// @Produce json
// @Param id path int true "Discipline ID"
// @Param payload body disciplineRequest true "Discipline payload"
// @Success 200 {object} models.Discipline
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /disciplines/{id} [put]
func (h *DisciplineHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req disciplineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	d := &models.Discipline{
		ID:          id,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		TeamSize:    req.TeamSize,
		Metadata:    req.Metadata,
	}
	if req.IsActive != nil {
		d.IsActive = *req.IsActive
	} else {
		d.IsActive = true
	}
	if len(d.Metadata) == 0 {
		d.Metadata = json.RawMessage(`{}`)
	}

	if err := h.svc.Update(c.Request.Context(), d); err != nil {
		if errors.Is(err, repository.ErrDisciplineNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondData(c, http.StatusOK, d, nil)
}

// @Summary Delete discipline
// @Tags Disciplines
// @Produce json
// @Param id path int true "Discipline ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /disciplines/{id} [delete]
func (h *DisciplineHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrDisciplineNotFound) {
			RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondData(c, http.StatusNoContent, nil, nil)
}
