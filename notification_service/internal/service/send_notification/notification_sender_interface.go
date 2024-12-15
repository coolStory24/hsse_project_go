package send_notification

import (
	"notification_service/internal/models"
)

type INotificationSender interface {
	Send(notification models.NotificationData) error
}
