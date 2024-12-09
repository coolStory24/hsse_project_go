package responses

import "github.com/google/uuid"

type GetHotelResponse struct {
	Id         uuid.UUID `json:"id"`
	HotelName  string    `json:"hotel_name"`
	NightPrice int       `json:"night_price"`
	AdminId    uuid.UUID `json:"admin_id"`
}
