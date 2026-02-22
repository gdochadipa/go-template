package main

import (
	"log/slog"
	"net"
	"os"

	"github.com/spf13/viper"
	userv1 "github.com/user/go-templates/template-grpc-ddd/gen/go/user/v1"
	handler "github.com/user/go-templates/template-grpc-ddd/internal/adapter/handler/grpc"
	"github.com/user/go-templates/template-grpc-ddd/internal/adapter/storage/memory"
	"github.com/user/go-templates/template-grpc-ddd/internal/core/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	// Dependency Injection
	// 1. Adapters (Storage)
	userRepo := memory.NewUserRepository()

	// 2. Core (Service)
	userSvc := service.NewUserService(userRepo, logger)

	// 3. Adapters (Handler)
	userHandler := handler.NewUserHandler(userSvc)

	// gRPC Server Setup
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	s := grpc.NewServer()

	// Register generated service
	userv1.RegisterUserServiceServer(s, userHandler)

	// Register reflection for debugging (grpcurl)
	reflection.Register(s)

	logger.Info("gRPC server starting", "port", port)
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
