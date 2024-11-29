package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"hotel_service/internal/server"
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

	server.NewServer(cfg.ServerConfig, cfg.HotelService)
}

func loadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("file .env in the root of the project not found")
	}
	return nil
}
