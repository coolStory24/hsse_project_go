package requests

import (
	"github.com/google/uuid"
	"time"
)

type RentFilter struct {
	ClientID   *uuid.UUID
	HotelierID *uuid.UUID
	HotelID    *uuid.UUID
	FromDate   *time.Time
	ToDate     *time.Time
}
