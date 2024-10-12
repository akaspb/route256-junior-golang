package server

import (
	"context"
	"errors"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetCustomerOrders(ctx context.Context, req *pb.GetCustomerOrdersRequest) (*pb.GetCustomerOrdersResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	lastN := uint(0)
	if req.GetLastCount() != nil {
		lastN = uint(req.GetLastCount().GetValue())
	}

	orders, err := s.service.GetCustomerOrders(
		ctx,
		int64ToIDType(req.GetCustomerId()),
		lastN,
	)
	if err != nil {
		var argumentErr *service.ArgumentError
		if errors.As(err, &argumentErr) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	ordersProto := make([]*pb.CustomerOrderInfo, len(orders))
	for i, order := range orders {
		customerOrderInfo := orderIDWithExpiryAndStatusToCustomerOrderInfo(order)
		ordersProto[i] = &customerOrderInfo
	}

	return &pb.GetCustomerOrdersResponse{
		Orders: ordersProto,
	}, nil
}
