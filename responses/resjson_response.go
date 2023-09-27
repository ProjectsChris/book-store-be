package responses

import (
	"github.com/gin-gonic/gin"
)

type ResponseJSON struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ResponseMessage function response a JSON of the operation
func ResponseMessage(c *gin.Context, status int, msg string, data any) {
	c.JSON(status, ResponseJSON{
		Status:  status,
		Message: msg,
		Data:    data,
	})
}
