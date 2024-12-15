package notification_service

import (
	"booking_service/internal/rest/dtos/responses"
	"booking_service/internal/service_interaction/user_service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type NotificationData struct {
	UserContactData *user_service.UserData     `json:"user_contact_data"`
	RentData        *responses.GetRentResponse `json:"rent_data"`
}

type INotificationServiceBridge interface {
	SendNotification(ctx context.Context, notificationData *NotificationData)
}

type NotificationServiceBridge struct {
	writer *kafka.Writer
	tracer trace.Tracer
}

func NewNotificationServiceBridge(broker string, topic string) *NotificationServiceBridge {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	tracer := otel.Tracer("notification_service_bridge")
	return &NotificationServiceBridge{writer: writer, tracer: tracer}
}

func (b *NotificationServiceBridge) SendNotification(ctx context.Context, notificationData *NotificationData) {
	ctx, span := b.tracer.Start(ctx, "SendNotification",
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", b.writer.Stats().Topic),
			attribute.String("messaging.operation", "send"),
		),
	)
	defer span.End()

	jsonData, err := json.Marshal(notificationData)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to serialize notification data: %v", err))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to serialize notification data")
		return
	}

	span.SetAttributes(attribute.Int("messaging.message.size", len(jsonData)))

	message := kafka.Message{
		Value: jsonData,
	}

	err = b.writer.WriteMessages(ctx, message)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to send Kafka message: %v", err))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send Kafka message")
		return
	}

	slog.Info("Message sent successfully to Kafka")
	span.SetStatus(codes.Ok, "Message sent successfully")
}
