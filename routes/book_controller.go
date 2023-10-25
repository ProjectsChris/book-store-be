package routes

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

// BookRoutes function with all routers of book
func BookRoutes(r *gin.RouterGroup, db *sql.DB, ctx context.Context) {
	sqlDb := new(DatabaseSql)
	sqlDb.Db = db

	v1 := r.Group("/api/v1")
	{
		book := v1.Group("/book")
		{
			// GET request
			book.GET("/", sqlDb.GetBooks)
			book.GET("/:title", sqlDb.GetBook)

			// POST request
			book.POST("/", sqlDb.PostBook)

			// PUT request
			book.PUT("/:id", sqlDb.UpdateBook)

			// DELETE request
			book.DELETE("/:id", sqlDb.DeleteBook)
		}
	}
}
