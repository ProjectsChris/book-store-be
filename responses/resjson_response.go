package responses

import (
	"book-store-be/models"
	"github.com/gin-gonic/gin"
)

// ResponseDatabase type for manage pagination of the database
type ResponseDatabase struct {
	Data       []models.Book      `json:"data"`
	Pagination PaginationDatabase `json:"paging"`
}

type PaginationDatabase struct {
	TotalRecord int `json:"total_record"`
	Page        int `json:"page"`
	TotalPages  int `json:"total_pages"`
}

// ResponseErrorJSON type used for every time there is an error
type ResponseErrorJSON struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ResponseMessage response a JSON with status error and message error
func ResponseMessage(c *gin.Context, status int, msg string) {
	c.JSON(status, ResponseErrorJSON{
		Status:  status,
		Message: msg,
	})
}
