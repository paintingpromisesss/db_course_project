package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"db_course_project/internal/pagination"
)

// ParsePagination extracts limit/offset from query params with defaults.
func ParsePagination(c *gin.Context) (limit, offset int) {
	limit = pagination.DefaultLimit
	offset = 0

	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if v := c.Query("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	limit, offset = pagination.Normalize(limit, offset)
	return
}
