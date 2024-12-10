package models

import (
	"github.com/google/uuid"
	"time"
)

type BookingRequest struct {
	FromDate    time.Time
	ToDate      time.Time
	HotelID     uuid.UUID
	ClientID    uuid.UUID
	ClientEmail string
}
