package server

import (
	"context"
	"errors"
	"fmt"
	"hotel_service/internal/config"
	"hotel_service/internal/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewServer(cfg *config.ServerConfig, hotelService services.IHotelService) {
	router := SetupApiRouter(cfg, hotelService)

	// Server configuration
	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	fmt.Printf("Server is starting on localhost%s\n", cfg.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Could not listen on %s: %v\n", cfg.Port, err)
		}
	}()

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
