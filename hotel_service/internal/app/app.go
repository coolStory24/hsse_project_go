package app

import (
	"hotel_service/internal/config"
	"hotel_service/internal/server"
)

func StartApp() {
	cfg, err := config.GetServerConfig()

	if err != nil {
		panic(err)
	}

	server.NewServer(cfg)
}
