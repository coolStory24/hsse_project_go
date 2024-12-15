package requests

import (
	"github.com/google/uuid"
	"time"
)

type CreateRentRequest struct {
	HotelID      uuid.UUID `json:"hotel_id"`
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
}
