package server

import (
	"booking_service/internal/config"
	db2 "booking_service/internal/db"
	"booking_service/internal/metrics"
	"booking_service/internal/service_interaction/hotel_service"
	"booking_service/internal/service_interaction/notification_service"
	"booking_service/internal/service_interaction/user_service"
	"booking_service/internal/services"
	"booking_service/internal/tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
	"os"
)

type CommonConfiguration struct {
	ServerConfig   *config.ServerConfig
	BookingService services.IBookingService
	TracerProvider *trace.TracerProvider
}

func NewCommonConfiguration() (*CommonConfiguration, error) {
	slog.Info("Creating common configuration")

	slog.Info("Getting configs")
	cfg, err := config.GetServerConfig()
	if err != nil {
		slog.Error("Failed to get configs")
		return nil, err
	}

	// Initialize tracing
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	tracerProvider, err := tracing.InitTracerProvider("booking_service", jaegerEndpoint)
	if err != nil {
		slog.Error("Failed to initialize tracing")
		return nil, err
	}

	// load database
	db, err := db2.NewDatabase()
	if err != nil {
		slog.Info("Failed to establish connection to database")
		return nil, err
	}
	slog.Info("Connection to database established")

	// setup grpc with hotel service
	hotelServiceBridge, err := hotel_service.NewHotelServiceBridge(os.Getenv("hotel_service_url"))
	if err != nil {
		slog.Info("Failed to establish gRPC connection with hotel service")
		return nil, err
	}
	slog.Info("gRPC connection with hotel service established")

	// setup grpc with user service
	userServiceBridge, err := user_service.NewUserServiceBridge(os.Getenv("user_service_url"))
	if err != nil {
		slog.Info("Failed to establish gRPC connection with user service")
		return nil, err
	}
	slog.Info("gRPC connection with user service established")

	// setup kafka to notification service
	notificationServiceKafkaBroker := os.Getenv("notification_service_kafka_broker")
	notificationServiceKafkaTopic := os.Getenv("notification_service_kafka_topic")
	notificationServiceBridge := notification_service.NewNotificationServiceBridge(
		notificationServiceKafkaBroker,
		notificationServiceKafkaTopic)
	slog.Info("Connection to kafka broker with notification service established")

	// setup metrics
	metrics.Register()
	slog.Info("Metrics registered")

	bookingService := services.NewBookingService(db, hotelServiceBridge, userServiceBridge, notificationServiceBridge)
	slog.Info("Booking service taken up")

	slog.Info("Common configuration was successfully created")
	return &CommonConfiguration{
		ServerConfig:   cfg,
		BookingService: bookingService,
		TracerProvider: tracerProvider,
	}, nil
}
