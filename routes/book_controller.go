package routes

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

var Tracer = otel.Tracer("Book Store Be")

// BookRoutes function with all routers of book
func BookRoutes(r *gin.RouterGroup, db *sql.DB, ctx context.Context) {
	//var bookSpan trace.Span

	sqlDb := new(DatabaseSql)
	sqlDb.Db = db

	v1 := r.Group("/api/v1")
	{
		book := v1.Group("/book")
		{
			// creates a span parent
			_, bookSpan := Tracer.Start(ctx, "/book")
			defer bookSpan.End()

			// GET request
			book.GET("/", sqlDb.GetBooks)
			book.GET("/:title", sqlDb.GetBook)

			// POST request
			book.POST("/", sqlDb.PostBook)
		}
	}

}
