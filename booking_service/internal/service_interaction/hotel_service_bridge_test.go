package service_interaction_test

import (
	sarama_mocks "booking_service/internal/mocks"
	"booking_service/internal/service_interaction"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestSendKafkaMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := sarama_mocks.NewMockSyncProducer(ctrl)

	bridge := service_interaction.NewHotelServiceBridge(mockProducer, nil)

	hotelID := uuid.New()

	mockProducer.EXPECT().
		SendMessage(gomock.Any()).
		DoAndReturn(func(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
			assert.Equal(t, service_interaction.KafkaTopicRequest, msg.Topic)
			assert.Equal(t, hotelID.String(), string(msg.Key.(sarama.StringEncoder)))
			assert.Equal(t, hotelID.String(), string(msg.Value.(sarama.StringEncoder)))
			return 0, 0, nil
		})

	err := bridge.SendKafkaMessage(hotelID)
	assert.NoError(t, err)
}

func TestReceiveKafkaMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumer := sarama_mocks.NewMockConsumer(ctrl)
	mockPartitionConsumer := sarama_mocks.NewMockPartitionConsumer(ctrl)

	bridge := service_interaction.NewHotelServiceBridge(nil, mockConsumer)

	hotelID := uuid.New()
	expectedPrice := 100
	responseMessage := struct {
		HotelID uuid.UUID `json:"hotel_id"`
		Price   int       `json:"price"`
	}{
		HotelID: hotelID,
		Price:   expectedPrice,
	}

	messageBytes, _ := json.Marshal(responseMessage)

	mockConsumer.EXPECT().
		ConsumePartition(service_interaction.KafkaTopicResponse, int32(0), sarama.OffsetNewest).
		Return(mockPartitionConsumer, nil)

	mockPartitionConsumer.EXPECT().
		Messages().
		DoAndReturn(func() <-chan *sarama.ConsumerMessage {
			ch := make(chan *sarama.ConsumerMessage, 1)
			go func() {
				ch <- &sarama.ConsumerMessage{Value: messageBytes}
			}()
			return ch
		})

	mockPartitionConsumer.EXPECT().Close()

	price, err := bridge.ReceiveKafkaMessage(hotelID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPrice, price)
}

func TestGetHotelPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := sarama_mocks.NewMockSyncProducer(ctrl)
	mockConsumer := sarama_mocks.NewMockConsumer(ctrl)
	mockPartitionConsumer := sarama_mocks.NewMockPartitionConsumer(ctrl)

	bridge := service_interaction.NewHotelServiceBridge(mockProducer, mockConsumer)

	hotelID := uuid.New()
	expectedPrice := 100
	responseMessage := struct {
		HotelID uuid.UUID `json:"hotel_id"`
		Price   int       `json:"price"`
	}{
		HotelID: hotelID,
		Price:   expectedPrice,
	}

	messageBytes, _ := json.Marshal(responseMessage)

	// Mock SendKafkaMessage
	mockProducer.EXPECT().
		SendMessage(gomock.Any()).
		DoAndReturn(func(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
			return int32(0), int64(0), nil
		})

	// Mock ConsumePartition
	mockConsumer.EXPECT().
		ConsumePartition(service_interaction.KafkaTopicResponse, int32(0), sarama.OffsetNewest).
		Return(mockPartitionConsumer, nil)

	// Mock Messages
	mockPartitionConsumer.EXPECT().
		Messages().
		DoAndReturn(func() <-chan *sarama.ConsumerMessage {
			ch := make(chan *sarama.ConsumerMessage, 1)
			go func() {
				ch <- &sarama.ConsumerMessage{Value: messageBytes}
			}()
			return ch
		})

	// Mock Close
	mockPartitionConsumer.EXPECT().Close()

	price, err := bridge.GetHotelPrice(hotelID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPrice, price)
}

func TestGetHotelPrice_SendKafkaMessageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := sarama_mocks.NewMockSyncProducer(ctrl)
	mockConsumer := sarama_mocks.NewMockConsumer(ctrl)

	bridge := service_interaction.NewHotelServiceBridge(mockProducer, mockConsumer)

	hotelID := uuid.New()

	// Mock SendKafkaMessage to return an error
	mockProducer.EXPECT().
		SendMessage(gomock.Any()).
		Return(int32(0), int64(0), fmt.Errorf("failed to send Kafka message"))

	price, err := bridge.GetHotelPrice(hotelID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, 0, price)
}

func TestReceiveKafkaMessageTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumer := sarama_mocks.NewMockConsumer(ctrl)
	mockPartitionConsumer := sarama_mocks.NewMockPartitionConsumer(ctrl)

	bridge := service_interaction.NewHotelServiceBridge(nil, mockConsumer)

	hotelID := uuid.New()

	mockConsumer.EXPECT().
		ConsumePartition(service_interaction.KafkaTopicResponse, int32(0), sarama.OffsetNewest).
		Return(mockPartitionConsumer, nil)

	// Mock Messages to return an empty channel
	mockPartitionConsumer.EXPECT().Messages().Return(make(chan *sarama.ConsumerMessage))
	mockPartitionConsumer.EXPECT().Close()

	// Ensure timeout behavior
	start := time.Now()
	price, err := bridge.ReceiveKafkaMessage(hotelID)
	assert.Error(t, err)
	assert.Equal(t, "timeout waiting for hotel price response", err.Error())
	assert.Equal(t, 0, price)
	assert.WithinDuration(t, start.Add(5*time.Second), time.Now(), 100*time.Millisecond)
}
