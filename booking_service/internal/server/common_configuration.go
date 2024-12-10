package server

import (
	"booking_service/internal/config"
	db2 "booking_service/internal/db"
	"booking_service/internal/metrics"
	"booking_service/internal/service_interaction/hotel_service"
	"booking_service/internal/service_interaction/notification_service"
	"booking_service/internal/service_interaction/user_service"
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

	// setup grpc with hotel service
	hotelServiceBridge, err := hotel_service.NewHotelServiceBridge(os.Getenv("hotel_service_url"))
	if err != nil {
		return nil, err
	}

	// setup grpc with user service
	userServiceBridge, err := user_service.NewUserServiceBridge(os.Getenv("user_service_url"))
	if err != nil {
		return nil, err
	}

	// setup kafka to notification service
	notificationServiceBridge, err := notification_service.NewNotificationServiceBridge()
	if err != nil {
		return nil, err
	}

	// setup metrics
	metrics.Register()

	bookingService := services.NewBookingService(db, hotelServiceBridge, userServiceBridge, notificationServiceBridge)

	return &CommonConfiguration{
		ServerConfig:   cfg,
		BookingService: bookingService,
	}, nil
}
