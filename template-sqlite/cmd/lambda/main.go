package main

import (
	"context"
	"database/sql"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/user/go-templates/template-sqlite/internal/config"
	"github.com/user/go-templates/template-sqlite/internal/user"
	"github.com/user/go-templates/template-sqlite/pkg/logger"
	_ "modernc.org/sqlite"
)

var chiLambda *chiadapter.ChiLambda

func init() {
	// Load Configuration
	cfg, err := config.LoadConfig("config")
	if err != nil {
		// handle error
	}

	// Initialize Logger
	log, _ := logger.New(cfg.Log.Level)

	// Connect to Database
	db, err := sql.Open(cfg.DB.Driver, cfg.DB.Source)
	if err != nil {
		// handle error
	}

	// Initialize Layers
	userRepo := user.NewSqliteRepository(db)
	userService := user.NewService(userRepo, log)
	userHandler := user.NewHandler(userService)

	// Router Setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)
	})

	chiLambda = chiadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
