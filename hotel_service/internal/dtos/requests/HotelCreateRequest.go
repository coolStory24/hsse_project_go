package requests

type CreateHotelRequest struct {
	HotelName  string `json:"hotel_name"`
	NightPrice int    `json:"night_price"`
}
