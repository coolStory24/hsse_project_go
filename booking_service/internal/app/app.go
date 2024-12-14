package app

import (
	"booking_service/internal/server"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func StartApp() {
	slog.Info("Launching the application")
	err := loadEnv()
	if err != nil {
		panic(err)
	}

	cfg, err := server.NewCommonConfiguration()

	if err != nil {
		panic(err)
	}

	server.NewServer(cfg.ServerConfig, cfg.BookingService)
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
		slog.Error(fmt.Sprintf("file %s not found in the root of the project: %w", fileName, err))
	}

	return nil
}
