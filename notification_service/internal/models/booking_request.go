package models

import (
	"github.com/google/uuid"
	"time"
)

type NotificationData struct {
	UserContactData *UserContactData `json:"user_contact_data"`
	RentData        *RentData        `json:"rent_data"`
}

type UserContactData struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type RentData struct {
	ID           uuid.UUID `json:"id"`
	HotelID      uuid.UUID `json:"hotel_id"`
	ClientID     uuid.UUID `json:"client_id"`
	NightPrice   int       `json:"night_price"`
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
}
