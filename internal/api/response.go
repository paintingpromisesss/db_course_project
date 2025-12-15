package api

import "github.com/gin-gonic/gin"

// PaginationMeta describes pagination info for list responses.
type PaginationMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// RespondData sends a standard JSON envelope.
func RespondData(c *gin.Context, status int, data any, meta any) {
	c.JSON(status, gin.H{
		"data": data,
		"meta": meta,
	})
}

// RespondError sends an error payload.
func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": message,
		},
	})
}
