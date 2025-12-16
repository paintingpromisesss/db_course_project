package api

import "github.com/gin-gonic/gin"

type PaginationMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func RespondData(c *gin.Context, status int, data any, meta any) {
	c.JSON(status, gin.H{
		"data": data,
		"meta": meta,
	})
}

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": message,
		},
	})
}
