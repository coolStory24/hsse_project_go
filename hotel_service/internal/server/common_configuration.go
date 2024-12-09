package server

import (
	"hotel_service/internal/config"
	"hotel_service/internal/metrics"
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

	// Create hotel service instance
	hotelService := &services.HotelService{}

	// Register metrics
	metrics.Register()

	return &CommonConfiguration{
		ServerConfig: cfg,
		HotelService: hotelService,
	}, nil
}
