package requests

type UpdateHotelRequest struct {
	HotelName  string `json:"hotel_name"`
	NightPrice int    `json:"night_price"`
}
