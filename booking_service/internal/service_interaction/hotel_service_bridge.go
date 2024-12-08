package service_interaction

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"time"
)

const (
	KafkaTopicRequest  = "hotel-request"
	KafkaTopicResponse = "hotel-response"
	KafkaBrokerAddress = "localhost:9092"
)

type IHotelServiceBridge interface {
	GetHotelPrice(hotelId uuid.UUID) (int, error)
	SendKafkaMessage(hotelId uuid.UUID) error
	ReceiveKafkaMessage(hotelId uuid.UUID) (int, error)
}

type HotelServiceBridge struct {
	kafkaProducer sarama.SyncProducer
	kafkaConsumer sarama.Consumer
}

func NewHotelServiceBridge(kafkaProducer sarama.SyncProducer, kafkaConsumer sarama.Consumer) *HotelServiceBridge {
	return &HotelServiceBridge{kafkaProducer, kafkaConsumer}
}

func CommonHotelServiceBridge() (*HotelServiceBridge, error) {
	bridge := &HotelServiceBridge{}

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

func (bridge *HotelServiceBridge) GetHotelPrice(hotelID uuid.UUID) (int, error) {
	err := bridge.SendKafkaMessage(hotelID)
	if err != nil {
		return 0, err
	}
	return bridge.ReceiveKafkaMessage(hotelID)
}

func (bridge *HotelServiceBridge) SendKafkaMessage(hotelID uuid.UUID) error {
	message := &sarama.ProducerMessage{
		Topic: KafkaTopicRequest,
		Key:   sarama.StringEncoder(hotelID.String()),
		Value: sarama.StringEncoder(hotelID.String()),
	}
	_, _, err := bridge.kafkaProducer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("error sending Kafka message: %w", err)
	}
	fmt.Println("Sent hotel price request for hotel ID:", hotelID)
	return nil
}

func (bridge *HotelServiceBridge) ReceiveKafkaMessage(hotelID uuid.UUID) (int, error) {
	partitionConsumer, err := bridge.kafkaConsumer.ConsumePartition(KafkaTopicResponse, 0, sarama.OffsetNewest)
	if err != nil {
		return 0, fmt.Errorf("error starting Kafka consumer: %w", err)
	}
	defer partitionConsumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var response struct {
				HotelID uuid.UUID `json:"hotel_id"`
				Price   int       `json:"price"`
			}
			err := json.Unmarshal(msg.Value, &response)
			if err != nil {
				return 0, fmt.Errorf("error decoding Kafka message: %w", err)
			}
			if response.HotelID == hotelID {
				return response.Price, nil
			}
		case <-ctx.Done():
			return 0, fmt.Errorf("timeout waiting for hotel price response")
		}
	}
}
