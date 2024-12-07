package service_interaction

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"hotel_service/internal/services"
)

const (
	KafkaTopicRequest  = "hotel-request"
	KafkaTopicResponse = "hotel-response"
	KafkaBrokerAddress = "localhost:9092"
)

type BookingServiceBridge struct {
	kafkaProducer sarama.SyncProducer
	kafkaConsumer sarama.Consumer
	hotelService  services.IHotelService
}

func CommonBookingServiceBridge() (*BookingServiceBridge, error) {
	bridge := &BookingServiceBridge{}

	var err error
	bridge.kafkaProducer, err = sarama.NewSyncProducer([]string{KafkaBrokerAddress}, nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing Kafka producer: %w", err)
	}

	bridge.kafkaConsumer, err = sarama.NewConsumer([]string{KafkaBrokerAddress}, nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing Kafka consumer: %w", err)
	}

	fmt.Println("Kafka producer and consumer initialized")
	return bridge, nil
}

func (bridge *BookingServiceBridge) StartListeningForHotelPriceRequests() {
	partitionConsumer, err := bridge.kafkaConsumer.ConsumePartition(KafkaTopicRequest, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Printf("Failed to start Kafka consumer: %v\n", err)
		return
	}
	defer partitionConsumer.Close()

	// Listen for messages
	for msg := range partitionConsumer.Messages() {
		// Decode the message
		var request struct {
			HotelID uuid.UUID `json:"hotel_id"`
		}
		err := json.Unmarshal(msg.Value, &request)
		if err != nil {
			fmt.Printf("Failed to decode Kafka message: %v\n", err)
			continue
		}

		hotel, err := bridge.hotelService.GetByID(request.HotelID)
		if err != nil {
			fmt.Printf("Failed to get price for hotel %s: %v\n", request.HotelID, err)
			continue
		}
		price := hotel.NightPrice

		// Prepare the response message
		response := struct {
			HotelID uuid.UUID `json:"hotel_id"`
			Price   int       `json:"price"`
		}{
			HotelID: request.HotelID,
			Price:   price,
		}

		// Send the response message back to Kafka
		responseMessage, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("Failed to encode Kafka response: %v\n", err)
			continue
		}

		message := &sarama.ProducerMessage{
			Topic: KafkaTopicResponse,
			Key:   sarama.StringEncoder(request.HotelID.String()),
			Value: sarama.StringEncoder(responseMessage),
		}

		_, _, err = bridge.kafkaProducer.SendMessage(message)
		if err != nil {
			fmt.Printf("Failed to send Kafka response: %v\n", err)
		} else {
			fmt.Printf("Sent hotel price response for hotel ID %s\n", request.HotelID)
		}
	}
}
