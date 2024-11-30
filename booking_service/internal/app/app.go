package app

import (
	"booking_service/internal/server"
	"fmt"
	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("file .env in the root of the project not found")
	}
	return nil
}
