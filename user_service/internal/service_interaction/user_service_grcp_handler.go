package service_interaction

import (
	"log/slog"
	pb "user_service/internal/service_interaction/gen"
	"user_service/internal/services"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingServiceBridge struct {
	pb.UnimplementedUserServiceServer
	userService services.IUserService
}

func NewUserServiceBridge(userService services.IUserService) *BookingServiceBridge {
	return &BookingServiceBridge{
		UnimplementedUserServiceServer: pb.UnimplementedUserServiceServer{},
		userService:                    userService,
	}
}

func (s *BookingServiceBridge) GetHotelPrice(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slog.Info("Handling request to get price of hotel with id " + req.Token)

	user, err := s.userService.GetUserByToken(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.GetUserResponse{Id: user.Id.String(), Email: user.Email, Username: user.Username, Role: user.Role.String()}, nil
}
