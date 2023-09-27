package routes

import (
	"book-store-be/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")
	{
		// GET request
		v1.GET("books", controllers.GetBooks)

		// POST request
		v1.POST("/new-book", controllers.PostBook)
	}
}
