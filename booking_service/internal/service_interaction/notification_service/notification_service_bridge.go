package notification_service

import (
	"booking_service/internal/rest/dtos/responses"
	"booking_service/internal/service_interaction/user_service"
)

type NotificationData struct {
	UserContactData *user_service.UserContactData
	RentData        *responses.GetRentResponse
}

type INotificationServiceBridge interface {
	SendNotification(notificationData *NotificationData)
}

type NotificationServiceBridge struct {
}

func NewNotificationServiceBridge() (*NotificationServiceBridge, error) {
	return &NotificationServiceBridge{}, nil
}

func (b *NotificationServiceBridge) SendNotification(notificationData *NotificationData) {
	// todo: send kafka request
}
