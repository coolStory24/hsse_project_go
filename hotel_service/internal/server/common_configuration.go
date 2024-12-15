package server

import (
	"go.opentelemetry.io/otel/sdk/trace"
	"hotel_service/internal/config"
	db2 "hotel_service/internal/db"
	"hotel_service/internal/metrics"
	"hotel_service/internal/services"
	"hotel_service/internal/tracing"
	"log/slog"
	"os"
)

type CommonConfiguration struct {
	ServerConfig   *config.ServerConfig
	HotelService   *services.HotelService
	TracerProvider *trace.TracerProvider
}

func NewCommonConfiguration() (*CommonConfiguration, error) {
	slog.Info("Creating common configuration")
	cfg, err := config.GetServerConfig()

	if err != nil {
		return nil, err
	}

	// Initialize tracing
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	tracerProvider, err := tracing.InitTracerProvider("hotel_service", jaegerEndpoint)
	if err != nil {
		slog.Error("Failed to initialize tracing")
		return nil, err
	}
	slog.Info("Created tracing provider")

	db, err := db2.NewDatabase()
	if err != nil {
		return nil, err
	}
	slog.Info("Connection to database established")

	hotelService := services.NewHotelService(db)

	// Register metrics
	metrics.Register()
	slog.Info("Metrics registered")

	slog.Info("Common configuration was successfully created")
	return &CommonConfiguration{
		ServerConfig:   cfg,
		HotelService:   hotelService,
		TracerProvider: tracerProvider,
	}, nil
}
