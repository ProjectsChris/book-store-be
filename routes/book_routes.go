package routes

import (
	"book-store-be/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")
	{
		// POST request
		v1.POST("/new-book", controllers.PostBook)
	}
}
