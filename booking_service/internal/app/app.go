package app

import (
	"booking_service/internal/server"
	"fmt"
	"github.com/joho/godotenv"
	"os"
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
		return fmt.Errorf("file %s not found in the root of the project: %w", fileName, err)
	}

    slog.Info("File .env was successfully loaded")
	return nil
}
