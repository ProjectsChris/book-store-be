package main

import (
	"book-store-be/database"
	_ "book-store-be/docs"
	"book-store-be/middleware"
	"book-store-be/observability"
	"book-store-be/routes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Book Store API
// @version		1.0
// @description	API for manages a book store
// @contact.name	Chris Developer
// @contact.email	chrisd3v3l0p3r@gmail.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			192.168.3.6:8000
// @BasePath		/api/v1
func main() {
	// creates a context
	ctx := context.Background()

	// read configuration file
	config, err := ReadConfig()
	if err != nil {
		log.Fatal("Impossibile Leggere il file di configurazione")
	}

	if config.Observability.Enable == false {
		// init tracer
		trace, err := observability.InitTracer(ctx, config.Observability.Endpoint, config.Observability.ServiceName)
		if err != nil {
			panic("trace error" + err.Error())
		}
		defer trace(ctx)

		// TODO: fix metrics
		// init metric
		//metric, err := observability.InitMetric(ctx, config.Observability.Endpoint, config.Observability.ServiceName)
		//if err != nil {
		//	panic("metric error" + err.Error())
		//}
		//defer metric(ctx)
	}

	// create a connection string for Postgresql
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.ConnectionStringPostgres.Host,
		config.Database.ConnectionStringPostgres.Port,
		config.Database.ConnectionStringPostgres.User,
		config.Database.ConnectionStringPostgres.Password,
		config.Database.ConnectionStringPostgres.DbName,
		config.Database.ConnectionStringPostgres.SslMode,
	)

	// connection to the database
	sqlDatabase := database.InitDatabase(connectionString)

	// gin
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// routes book
	routes.BookRoutes(&r.RouterGroup, sqlDatabase, ctx)

	// swagger API
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Status":      "UP!",
			"Description": "HealthCheck",
		})
	})

	// run gin
	r.Run("0.0.0.0:8000")
}
