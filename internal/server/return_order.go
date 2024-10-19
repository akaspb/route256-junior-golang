package server

import (
	"context"
	"errors"

	pb "gitlab.ozon.dev/siralexpeter/Homework/internal/pvz-service/v1"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) ReturnOrder(ctx context.Context, req *pb.ReturnOrderRequest) (*pb.ReturnOrderResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.ReturnOrderFromCustomer(
		ctx,
		int64ToIDType(req.GetCustomerId()),
		int64ToIDType(req.GetOrderId()),
	)
	if err != nil {
		var argumentErr *service.ArgumentError
		if errors.As(err, &argumentErr) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ReturnOrderResponse{}, nil
}
