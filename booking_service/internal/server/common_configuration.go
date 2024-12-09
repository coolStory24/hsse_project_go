package server

import (
	"booking_service/internal/config"
	db2 "booking_service/internal/db"
	"booking_service/internal/metrics"
	"booking_service/internal/service_interaction"
	"booking_service/internal/services"
	"os"
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

	// load database
	db, err := db2.NewDatabase()
	if err != nil {
		return nil, err
	}

	// setup grpc
	bridge, err := service_interaction.NewHotelServiceBridge(os.Getenv("hotel_service_url"))
	if err != nil {
		return nil, err
	}

	// setup metrics
	metrics.Register()

	bookingService := services.NewBookingService(db, bridge)

	return &CommonConfiguration{
		ServerConfig:   cfg,
		BookingService: bookingService,
	}, nil
}
