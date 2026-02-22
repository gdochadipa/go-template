package main

import (
	"log"
	"net"

	userv1 "github.com/user/go-templates/template-grpc-sdk/gen/go/user/v1"
	"github.com/user/go-templates/template-grpc-sdk/internal/config"
	"github.com/user/go-templates/template-grpc-sdk/internal/user"
	"github.com/user/go-templates/template-grpc-sdk/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := logger.New(cfg.Log.Level)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register services
	userSvc := user.NewService(logger)
	userv1.RegisterUserServiceServer(s, userSvc)

	// Register reflection service on gRPC server.
	reflection.Register(s)

	logger.Info("gRPC server starting", zap.String("port", cfg.Server.Port))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
