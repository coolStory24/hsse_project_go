package hotel_service

import (
	gen2 "booking_service/internal/service_interaction/hotel_service/gen"
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type IHotelServiceBridge interface {
	GetHotelPrice(hotelId uuid.UUID) (int, error)
}

type HotelServiceBridge struct {
	grpcClient gen2.HotelServiceClient
}

func NewHotelServiceBridge(grpcAddress string) (*HotelServiceBridge, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := gen2.NewHotelServiceClient(conn)

	return &HotelServiceBridge{grpcClient: client}, nil
}

func (h *HotelServiceBridge) GetHotelPrice(hotelId uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &gen2.GetHotelPriceRequest{HotelId: hotelId.String()}
	slog.Info("Sending request to get price of hotel with id " + hotelId.String())
	response, err := h.grpcClient.GetHotelPrice(ctx, request)
	if err != nil {
		return 0, err
	}
	return int(response.Price), nil
}
