package main

import (
	"book-store-be/database"
	"book-store-be/observability"
	"book-store-be/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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

// @host		192.168.3.6:8000
// @BasePath	/api/v1
func main() {
	// init a context
	ctx := context.Background()

	// read configuration file
	config, err := ReadConfig()
	if err != nil {
		log.Fatal("Impossibile Leggere il file di configurazione")
	}

	// init tracer
	trace, err := observability.InitTracer(ctx, config.Observability.Endpoint)
	if err != nil {
		panic("trace error" + err.Error())
	}
	defer trace(ctx)

	// init metric
	metric, err := observability.InitMetric(ctx, config.Observability.Endpoint)
	if err != nil {
		panic("metric error" + err.Error())
	}
	defer metric(ctx)

	// create a connection string
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

	// CORS configuration
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// CORS middleware
	r.Use(cors.New(corsConfig))

	// routes
	routes.BookRoutes(&r.RouterGroup, sqlDatabase)

	// swagger API
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run gin
	r.Run(":8000")
}
