package server

import (
	"hotel_service/internal/config"
	db2 "hotel_service/internal/db"
	"hotel_service/internal/services"
	"log/slog"
)

type CommonConfiguration struct {
	ServerConfig *config.ServerConfig
	HotelService *services.HotelService
}

func NewCommonConfiguration() (*CommonConfiguration, error) {
	slog.Info("Creating common configuration")
	cfg, err := config.GetServerConfig()

	if err != nil {
		return nil, err
	}

	db, err := db2.NewDatabase()
	if err != nil {
		return nil, err
	}

	hotelService := services.NewHotelService(db)

	slog.Info("Common configuration was successfully created")
	return &CommonConfiguration{
		ServerConfig: cfg,
		HotelService: hotelService,
	}, nil
}
