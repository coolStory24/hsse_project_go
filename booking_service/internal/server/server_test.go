package server

import (
	"booking_service/internal/config"
	"booking_service/internal/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupApiRouter(t *testing.T) {
	serverConfig := &config.ServerConfig{}
	bookingService := &services.BookingService{}
	router := SetupApiRouter(serverConfig, bookingService)

	assert.NotNil(t, router)
}
