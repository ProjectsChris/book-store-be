package routes

import (
	"book-store-be/controllers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func BookRoutes(r *gin.RouterGroup, db *sql.DB) {
	sqlDb := new(controllers.DatabaseSql)
	sqlDb.Db = db

	v1 := r.Group("/api/v1")
	{
		book := v1.Group("/book")
		{
			// GET request
			book.GET("/:title", sqlDb.GetBook)

			// POST request
			book.POST("/", sqlDb.PostBook)
		}
	}
}
