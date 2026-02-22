package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/user/go-templates/template-mongo/internal/config"
	"github.com/user/go-templates/template-mongo/internal/user"
	"github.com/user/go-templates/template-mongo/pkg/logger"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	client, err := mongoDriver.Connect(context.Background(), options.Client().ApplyURI(cfg.DB.URI))
	if err != nil {
		// handle error
	}
	db := client.Database(cfg.DB.Database)

	// Initialize Layers
	userRepo := user.NewMongoRepository(db)
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
