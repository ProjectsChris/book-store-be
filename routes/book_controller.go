package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func BookRoutes(r *gin.RouterGroup, db *sql.DB) {
	sqlDb := new(DatabaseSql)
	sqlDb.Db = db

	v1 := r.Group("/api/v1")
	{
		book := v1.Group("/book")
		{
			// GET request
			book.GET("/:title", sqlDb.GetBook)
			book.GET("/", sqlDb.GetBooks)

			// POST request
			book.POST("/", sqlDb.PostBook)
		}
	}
}
