package responses

import (
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
