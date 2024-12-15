package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"notification_service/internal/server"
	"notification_service/internal/service/content_build"
	"notification_service/internal/service/send_notification"
	"os"
	"strconv"
)

func StartApp() {
	err := loadEnv()
	if err != nil {
		slog.Error("Failed to load env")
	}

	// Setup email sender
	emailContentBuilder := content_build.NewEmailContentBuilder()
	smtpServer := os.Getenv("SMTP_SERVER")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		slog.Error("Failed to parse SMTP_SERVER")
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	emailSender := send_notification.NewEmailSender(emailContentBuilder, smtpServer, port, username, password)

	// Gather all senders (we only have email sending so far)
	senders := []send_notification.INotificationSender{emailSender}

	// Setup notification request handler
	notificationRequestHandler := server.NewNotificationRequestHandler(senders)

	// Setup Kafka consumer
	kafkaConsumerBroker := os.Getenv("KAFKA_CONSUMER_BROKER")
	kafkaConsumerTopic := os.Getenv("KAFKA_CONSUMER_TOPIC")
	server.StartKafkaConsumer(kafkaConsumerBroker, kafkaConsumerTopic, notificationRequestHandler)
}

func loadEnv() error {
	env := os.Getenv("GO_ENV")
	var fileName string

	if env == "" {
		fileName = ".env"
	} else if env == "dev" {
		fileName = ".env.dev"
	}

	err := godotenv.Load(fileName)
	if err != nil {
		return fmt.Errorf("file %s not found in the root of the project: %w", fileName, err)
	}

	return nil
}
