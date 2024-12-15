package content_build

import (
	"fmt"
	"log/slog"
	"notification_service/internal/models"
)

type EmailContentBuilder struct{}

func NewEmailContentBuilder() *EmailContentBuilder {
	return &EmailContentBuilder{}
}

func (e *EmailContentBuilder) BuildContent(notification models.NotificationData) string {
	booking := notification

	fromDate := booking.RentData.CheckInDate.Format("January 2, 2006")
	toDate := booking.RentData.CheckOutDate.Format("January 2, 2006")

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
		booking.RentData.HotelID, fromDate, toDate, booking.RentData.ClientID,
	)
	slog.Info("Built content of the email")

	return emailContent
}
