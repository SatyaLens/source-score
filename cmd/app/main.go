package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	embed "source-score"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"source-score/pkg/api"
	"source-score/pkg/conf"
	"source-score/pkg/db/pgsql"
	"source-score/pkg/domain/claim"
	"source-score/pkg/domain/proof"
	"source-score/pkg/domain/source"
	"source-score/pkg/handlers"
	"source-score/pkg/helpers"
	apiServer "source-score/pkg/http"
	"source-score/pkg/logger"
	"source-score/pkg/middleware"
)

func main() {
	// initialize the config
	conf.LoadConfig()

	// initialize the logger
	slog.SetDefault(
		slog.New(&logger.ContextHandler{
			Handler: slog.NewJSONHandler(os.Stdout, nil),
		}),
	)

	// initialize the server
	loggerOpts := api.GinServerOptions{
		Middlewares: []api.MiddlewareFunc{
			// function to add request headers to log fields
			func(c *gin.Context) {
				for _, fieldKey := range helpers.ApiReqLogFields {
					fieldVal := c.Request.Header.Get(fieldKey)
					logger.AppendGinCtx(c, slog.String(fieldKey, fieldVal))
				}
			},
		},
	}

	// initialize the layers
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		conf.Cfg.PgHost,
		conf.AppUserName,
		conf.Cfg.AppUserPassword,
		conf.DbName,
	)
	dbClient := pgsql.NewClient(context.Background(), dsn, &gorm.Config{})
	// TODO: wrap this and call it securely
	dbClient.SetAutoMigration(
		context.Background(),
		api.Source{},
		api.Claim{},
	)

	proofRepo := proof.NewProofRepository(context.Background(), dbClient)
	proofSvc := proof.NewProofService(context.Background(), proofRepo)

	claimRepo := claim.NewClaimRepository(context.Background(), dbClient)
	claimSvc := claim.NewClaimService(context.Background(), claimRepo, proofSvc)

	srcRepo := source.NewSourceRepository(context.Background(), dbClient)
	srcSvc := source.NewSourceService(context.Background(), srcRepo, claimRepo)

	server := gin.Default()
	// register middlewares
	endpoints := server.Group("/api")
	// secure with API key if the env var is set
	if key, ok := os.LookupEnv("API_KEY"); ok {
		endpoints.Use(middleware.APIKeyMiddleware(key))
	}

	api.RegisterHandlersWithOptions(
		server,
		apiServer.NewRouter(context.Background(), srcSvc, claimSvc, proofSvc),
		loggerOpts,
	)

	// Register Swagger UI routes
	swaggerHandler := handlers.NewSwaggerHandler(embed.OpenAPI)
	server.GET("/swagger", swaggerHandler.ServeUI)
	server.GET("/swagger/spec", swaggerHandler.ServeSpec)

	err := server.Run()
	if err != nil {
		log.Fatalf("failed to start the server : %s\n", err.Error())
	} else {
		log.Println("Server is up and running!")
	}
}
