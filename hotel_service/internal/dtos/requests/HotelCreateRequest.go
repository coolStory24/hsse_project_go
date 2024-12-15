package requests

import "github.com/google/uuid"

type CreateHotelRequest struct {
	HotelName  string    `json:"hotel_name"`
	NightPrice int       `json:"night_price"`
	AdminId    uuid.UUID `json:"admin_id"`
}
