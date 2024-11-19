package app

import (
	"booking_service/internal/config"
	"booking_service/internal/server"
)

func StartApp() {
	cfg, err := config.GetServerConfig()

	if err != nil {
		panic(err)
	}

	server.NewServer(cfg)
}
