package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user_service/internal/config"
	"user_service/internal/db"
	"user_service/internal/repositories"
	"user_service/internal/rest"
	"user_service/internal/services"
)

func NewServer(cfg *config.ServerConfig) {
	dbConnection, err := db.NewDbConnection()

	if err != nil {
		fmt.Printf("Could not connect to the database %v\n", err)
		return
	}

	// migrationErr := migrator.RunMigrations(dbConnection.Connection)

	// if migrationErr != nil {
	// 	fmt.Println(err)
	// }

	router := rest.SetupApiRouter(cfg, services.NewUserService(repositories.NewUserRepository(dbConnection)))

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
