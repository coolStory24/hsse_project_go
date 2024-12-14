package send_notification

import (
	"notification_service/internal/models"
	"notification_service/internal/server2/content_build"
)

type INotificationSender interface {
	Send(notification models.NotificationData, contentBuilder content_build.IContentBuilder, clientData models.ClientData) error
}
