package server

import (
	"booking_service/internal/config"
	db2 "booking_service/internal/db"
	"booking_service/internal/services"
)

type CommonConfiguration struct {
	ServerConfig   *config.ServerConfig
	BookingService services.IBookingService
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

	bookingService := services.NewBookingService(db)

	return &CommonConfiguration{
		ServerConfig:   cfg,
		BookingService: bookingService,
	}, nil
}
