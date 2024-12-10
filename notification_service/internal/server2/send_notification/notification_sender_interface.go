package send_notification

import (
	"hotel_service/internal/models"
	"hotel_service/internal/server2/content_build"
)

type INotificationSender interface {
	Send(booking models.BookingRequest, contentBuilder content_build.IContentBuilder, clientData models.ClientData) error
}
