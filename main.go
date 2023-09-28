package main

import (
	"book-store-be/database"
	"book-store-be/routes"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "book-store-be/docs"

	"github.com/gin-gonic/gin"
)

//	@title			Book Store API
//	@version		1.0
//	@description	This API manage a cart

//	@contact.name	Chris Developer
//	@contact.email	chrisd3v3l0p3r@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		192.168.3.8:8080
// @BasePath	/api/v1
func main() {
	// connect to the database
	database.ConnectDatabase()

	// gin
	r := gin.Default()

	// routes
	routes.BookRoutes(&r.RouterGroup)

	// swagger API
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run gin
	r.Run()
}
