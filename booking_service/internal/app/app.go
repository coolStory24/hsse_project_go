package app

import (
	"booking_service/internal/server"
	"booking_service/internal/tracing"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func StartApp() {
	slog.Info("Launching the application")
	err := loadEnv()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	cfg, err := server.NewCommonConfiguration()

	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	defer tracing.ShutdownTracerProvider(context.Background(), cfg.TracerProvider)

	server.NewServer(cfg.ServerConfig, cfg.BookingService)
	slog.Info("Application is running")
}

func loadEnv() error {
	slog.Info("Loading .env file")
	env := os.Getenv("GO_ENV")
	var fileName string

	if env == "" {
		fileName = ".env"
	} else if env == "dev" {
		fileName = ".env.dev"
	}

	err := godotenv.Load(fileName)
	if err != nil {
		slog.Error(fmt.Sprintf("file %s not found in the root of the project: %e", fileName, err))
	}

	slog.Info("File .env was successfully loaded")
	return nil
}
