package server

import (
	"log/slog"
	"notification_service/internal/models"
	"notification_service/internal/service/send_notification"
)

type INotificationRequestHandler interface {
	HandleNotificationRequest(data models.NotificationData)
}

type NotificationRequestHandler struct {
	senders []send_notification.INotificationSender
}

func NewNotificationRequestHandler(senders []send_notification.INotificationSender) *NotificationRequestHandler {
	return &NotificationRequestHandler{
		senders: senders,
	}
}

func (h *NotificationRequestHandler) HandleNotificationRequest(data models.NotificationData) {
	for _, sender := range h.senders {
		err := sender.Send(data)
		if err != nil {
			slog.Error("Failed to send message")
		} else {
			slog.Info("Message send successfully")
		}
	}
}
