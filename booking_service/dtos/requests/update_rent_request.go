package requests

import "time"

type UpdateRentRequest struct {
	HotelID      string    `json:"hotel_id,omitempty"`
	ClientID     string    `json:"client_id,omitempty"`
	CheckInDate  time.Time `json:"check_in_date,omitempty"`
	CheckOutDate time.Time `json:"check_out_date,omitempty"`
}
