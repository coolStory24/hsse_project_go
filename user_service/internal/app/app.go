package app

import (
	"user_service/internal/config"
	"user_service/internal/server"
)

func StartApp() {
	cfg, err := config.GetServerConfig()

	if err != nil {
		panic(err)
	}

	server.NewServer(cfg)
}
