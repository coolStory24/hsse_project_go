package server

import (
	"context"
	"fmt"
	"hotel_service/internal/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(cfg *config.ServerConfig) {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix(cfg.Prefix).Subrouter()

	// Handle requests here:
	//apiRouter.HandleFunc("/{path}", rest.{handler_name}).Methods("{METHOD}")
	fmt.Println(apiRouter)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	fmt.Printf("Server is starting on localhost%s\n", cfg.Port)

	go func() {
		err := srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
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
		fmt.Printf("Server forced to shutdown: %v \n", err)
	}

	fmt.Println("Server exited gracefully")
}
