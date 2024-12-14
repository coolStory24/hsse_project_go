package notification_service

import (
	"booking_service/internal/rest/dtos/responses"
	"booking_service/internal/service_interaction/user_service"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type NotificationData struct {
	UserContactData *user_service.UserContactData `json:"user_contact_data"`
	RentData        *responses.GetRentResponse    `json:"rent_data"`
}

type INotificationServiceBridge interface {
	SendNotification(notificationData *NotificationData)
}

type NotificationServiceBridge struct {
	writer *kafka.Writer
}

func NewNotificationServiceBridge(broker string, topic string) *NotificationServiceBridge {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &NotificationServiceBridge{writer: writer}
}

func (b *NotificationServiceBridge) SendNotification(notificationData *NotificationData) {
	jsonData, err := json.Marshal(notificationData)
	if err != nil {
		slog.Error("Failed to serialize notification data: %v", err)
		return
	}

	message := kafka.Message{
		Value: jsonData,
	}

	err = b.writer.WriteMessages(context.Background(), message)
	if err != nil {
		slog.Error("Failed to send Kafka message: %v", err)
		return
	}

	slog.Info("Message sent successfully to Kafka")
}
