package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user_service/internal/config"
	"user_service/internal/db"
	"user_service/internal/repositories"
	"user_service/internal/rest"
	"user_service/internal/service_interaction"
	pb "user_service/internal/service_interaction/gen"
	"user_service/internal/services"

	"google.golang.org/grpc"
)

func NewServer(cfg *config.ServerConfig) {
	dbConnection, err := db.NewDbConnection()

	if err != nil {
		fmt.Printf("Could not connect to the database %v\n", err)
		return
	}

	userService := services.NewUserService(repositories.NewUserRepository(dbConnection), services.NewEncryptionService(cfg.EncryptionKey))

	router := rest.SetupApiRouter(cfg, userService)

	// Server configuration
	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	fmt.Printf("Server is starting on port %s\n", cfg.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Could not listen on %s: %v\n", cfg.Port, err)
		}
	}()

	// gRCP
	grpcHotelService := service_interaction.NewUserServiceBridge(userService)
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpcHotelService)

	grpcListener, err := net.Listen("tcp", ":"+os.Getenv("GRCP_PORT"))
	if err != nil {
		slog.Error("Failed to listen on gRPC port")
		return
	}

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			slog.Error(fmt.Sprintf("gRPC server failed to start: %v\n", err))
		}
	}()

	slog.Info("Server is starting on localhost" + srv.Addr)
	slog.Info("gRPC server is starting on: " + os.Getenv("GRCP_PORT"))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server exited gracefully")
}
