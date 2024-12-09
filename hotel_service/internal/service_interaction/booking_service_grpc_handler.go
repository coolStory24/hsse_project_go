package service_interaction

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "hotel_service/internal/service_interaction/gen"
	"hotel_service/internal/services"
	"log/slog"
)

type BookingServiceBridge struct {
	pb.UnimplementedHotelServiceServer
	hotelService services.IHotelService
}

func NewBookingServiceBridge(hotelService services.IHotelService) *BookingServiceBridge {
	return &BookingServiceBridge{
		UnimplementedHotelServiceServer: pb.UnimplementedHotelServiceServer{},
		hotelService:                    hotelService,
	}
}

func (s *BookingServiceBridge) GetHotelPrice(ctx context.Context, req *pb.GetHotelPriceRequest) (*pb.GetHotelPriceResponse, error) {
	slog.Info("Handling request to get price of hotel with id " + req.HotelId)

	id, err := uuid.Parse(req.HotelId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid hotel ID: %v", err)
	}

	hotel, err := s.hotelService.GetByID(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "hotel not found: %v", err)
	}

	return &pb.GetHotelPriceResponse{Price: int32(hotel.NightPrice)}, nil
}
