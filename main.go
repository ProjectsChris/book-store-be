package main

import (
	"book-store-be/database"
	"book-store-be/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// connect to the databasa
	database.ConnectDatabase()

	// gin
	r := gin.Default()

	routes.BookRoutes(&r.RouterGroup)

	r.Run()
}
