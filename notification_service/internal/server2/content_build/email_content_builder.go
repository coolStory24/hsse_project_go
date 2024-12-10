package content_build

import (
	"fmt"
	"hotel_service/internal/models"
)

type EmailContentBuilder struct{}

func (e *EmailContentBuilder) BuildContent(bookingRequest models.BookingRequest) string {
	booking := bookingRequest

	fromDate := booking.FromDate.Format("January 2, 2006")
	toDate := booking.ToDate.Format("January 2, 2006")

	emailContent := fmt.Sprintf(
		"Dear Customer,\n\n"+
			"Thank you for your booking at our hotel. Below are your booking details:\n\n"+
			"Hotel ID: %s\n"+
			"Check-in Date: %s\n"+
			"Check-out Date: %s\n"+
			"Client Email: %s\n\n"+
			"We look forward to welcoming you. If you have any questions or need further assistance, feel free to reach out.\n\n"+
			"Best regards,\n"+
			"Your Hotel Team",
		booking.HotelID, fromDate, toDate, booking.ClientEmail,
	)

	return emailContent
}
