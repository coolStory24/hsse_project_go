package responses

type GetHotelsResponse struct {
	Hotels []GetHotelResponse `json:"hotels"`
}
