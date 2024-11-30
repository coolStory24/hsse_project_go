package responses

import (
	"github.com/google/uuid"
	"time"
)

type GetRentResponse struct {
	ID           uuid.UUID `json:"id"`
	HotelID      uuid.UUID `json:"hotel_id"`
	HotelierId   uuid.UUID `json:"hotelier_id"`
	ClientID     uuid.UUID `json:"client_id"`
	NightPrice   int       `json:"night_price"`
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
}
