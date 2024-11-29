package app

import (
	"hotel_service/internal/server"
)

func StartApp() {
	cfg, err := server.NewCommonConfiguration()

	if err != nil {
		panic(err)
	}

	server.NewServer(cfg.ServerConfig, cfg.HotelService)
}
