package server

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"notification_service/internal/models"
)

func StartKafkaConsumer(broker string, topic string, handler INotificationRequestHandler) {
	reader := createConsumer(broker, topic)
	defer reader.Close()

	log.Println("Consumer started. Waiting for messages...")
	StartConsumeLoop(reader, handler)
}

func createConsumer(broker string, topic string) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "notification_consumer_group",
	})
	return reader
}

func StartConsumeLoop(reader *kafka.Reader, handler INotificationRequestHandler) {
	ctx := context.Background()
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		var notificationData models.NotificationData
		err = json.Unmarshal(msg.Value, &notificationData)
		if err != nil {
			log.Printf("Failed to decode message: %v", err)
			continue
		}

		handler.HandleNotificationRequest(notificationData)
	}
}
