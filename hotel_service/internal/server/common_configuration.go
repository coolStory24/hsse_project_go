package server

import (
	"hotel_service/internal/config"
	db2 "hotel_service/internal/db"
	"hotel_service/internal/services"
)

type CommonConfiguration struct {
	ServerConfig *config.ServerConfig
	HotelService *services.HotelService
}

func NewCommonConfiguration() (*CommonConfiguration, error) {
	cfg, err := config.GetServerConfig()

	if err != nil {
		return nil, err
	}

	db, err := db2.NewDatabase()
	if err != nil {
		return nil, err
	}

	hotelService := services.NewHotelService(db)

	return &CommonConfiguration{
		ServerConfig: cfg,
		HotelService: hotelService,
	}, nil
}
