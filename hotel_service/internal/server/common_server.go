package server

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"hotel_service/internal/config"
	"hotel_service/internal/service_interaction"
	pb "hotel_service/internal/service_interaction/gen"
	"hotel_service/internal/services"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewServer(cfg *config.ServerConfig, hotelService services.IHotelService) {
	slog.Info("Starting a server")
	router := SetupApiRouter(cfg, hotelService)

	// Server configuration
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	slog.Info("Server is starting on localhost" + srv.Addr)
	// gRPC Server setup
	grpcHotelService := service_interaction.NewBookingServiceBridge(hotelService)
	grpcServer := grpc.NewServer()
	pb.RegisterHotelServiceServer(grpcServer, grpcHotelService)

	// Listener for gRPC
	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Failed to listen on gRPC port %s: %v\n", ":50051", err)
		return
	}

	fmt.Printf("Server is starting on localhost%s\n", srv.Addr)
	fmt.Printf("gRPC server is starting on %s\n", ":50051")

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Could not listen on %s: %v\n", cfg.Port, err)
		}
	}()

	// Start gRPC Server in a goroutine
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			fmt.Printf("gRPC server failed to start: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown: %v\n", err)
	}

	slog.Info("Server exited gracefully")
}
