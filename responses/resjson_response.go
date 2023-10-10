package responses

import (
	"book-store-be/models"
	"github.com/gin-gonic/gin"
)

type ResponseErrorJSON struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ResponseMessage function response a JSON of the operation
func ResponseMessage(c *gin.Context, status int, msg string) {
	c.JSON(status, ResponseErrorJSON{
		Status:  status,
		Message: msg,
	})
}

type ResponseDatabase struct {
	Data       []models.Book      `json:"data"`
	Pagination PaginationDatabase `json:"paging"`
}

type PaginationDatabase struct {
	TotalRecord int `json:"total_record"`
	Page        int `json:"page"`
	TotalPages  int `json:"total_pages"`
}
