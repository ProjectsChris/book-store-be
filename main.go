package main

import (
	"book-store-be/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// connect to the databasa
	database.ConnectDatabase()

	// gin
	r := gin.Default()
	r.Run()
}
