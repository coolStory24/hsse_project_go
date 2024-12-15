package content_build

import "notification_service/internal/models"

type IContentBuilder interface {
	BuildContent(notification models.NotificationData) string
}
