package responses

type GetRentsResponse struct {
	Rents []GetRentResponse `json:"rents"`
}
