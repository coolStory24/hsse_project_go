// This configuration can be ignored from test coverage, because it only injects the dependencies
// of services implementations and does not perform any business-logic
//go:build testnocover

package server

import (
	"hotel_service/internal/config"
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

	hotelService := &services.HotelService{}

	return &CommonConfiguration{
		ServerConfig: cfg,
		HotelService: hotelService,
	}, nil
}
