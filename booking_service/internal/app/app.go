package app

import (
	"booking_service/internal/server"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func StartApp() {
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

	return nil
}
