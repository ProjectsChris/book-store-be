package routes

import (
	"book-store-be/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func BookRoutes(r *gin.RouterGroup, db *mongo.Client) {

	mongoClient := new(controllers.MongoClient)
	mongoClient.Db = db 

	v1 := r.Group("/api/v1")
	{
		book := v1.Group("/book")
		{
			// GET request
			book.GET("/:title", mongoClient.GetBook)

			// POST request
			book.POST("/", mongoClient.PostBook)
		}
	}
}
