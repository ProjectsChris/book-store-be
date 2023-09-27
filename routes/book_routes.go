package routes

import (
	"book-store-be/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")
	{
		book := v1.Group("/book")
		{
			// GET request
			book.GET("/:titolo", controllers.GetBook)

			// POST request
			book.POST("/", controllers.PostBook)
		}
	}
}
