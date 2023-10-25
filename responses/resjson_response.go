package responses

import (
	"book-store-be/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponsePagination type for manage pagination of the database
type ResponsePagination struct {
	Data               []models.Book   `json:"data"`
	PaginationDatabase `json:"paging"` // struct embedding
}

type PaginationDatabase struct {
	TotalRecord int `json:"total_record"`
	Page        int `json:"page"`
	TotalPages  int `json:"total_pages"`
}

type Response struct {
	Message string `json:"message"`
}

func ErrorServerResponseJson(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Message: msg,
	})
}
