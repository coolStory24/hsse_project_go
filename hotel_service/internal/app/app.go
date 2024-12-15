package app

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"hotel_service/internal/server"
	"hotel_service/internal/tracing"
	"log/slog"
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

	server.NewServer(cfg.ServerConfig, cfg.HotelService)
}

func loadEnv() error {
	slog.Info("Loading .env file")
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("file .env in the root of the project not found")
	}
	slog.Info("File .env was successfully loaded")
	return nil
}
