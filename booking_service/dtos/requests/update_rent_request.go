package requests

import (
	"github.com/google/uuid"
	"time"
)

type UpdateRentRequest struct {
	HotelID      uuid.UUID `json:"hotel_id,omitempty"`
	ClientID     uuid.UUID `json:"client_id,omitempty"`
	CheckInDate  time.Time `json:"check_in_date,omitempty"`
	CheckOutDate time.Time `json:"check_out_date,omitempty"`
}
