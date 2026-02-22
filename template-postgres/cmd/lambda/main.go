package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/go-templates/template-postgres/internal/config"
	"github.com/user/go-templates/template-postgres/internal/user"
	"github.com/user/go-templates/template-postgres/pkg/logger"
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
	// In Lambda, connection pooling needs care.
	// We'll initialize it here (Cold Start).
	dbPool, err := pgxpool.New(context.Background(), cfg.DB.Source)
	if err != nil {
		// handle error
	}
	// Note: We don't defer close here because init runs once per cold start.
	// The connection stays open for warm invocations.

	// Initialize Layers
	userRepo := user.NewPostgresRepository(dbPool)
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
