package server

import (
	"context"
	"errors"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetReturnedOrders(ctx context.Context, req *pb.GetReturnedOrdersRequest) (*pb.GetReturnedOrdersResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	returns, err := s.service.GetReturnsList(ctx, int(req.GetOffset()), int(req.GetLimit()))

	if err != nil {
		var argumentErr *service.ArgumentError
		if errors.As(err, &argumentErr) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	ordersProto := make([]*pb.ReturnedOrder, len(returns))
	for i, order := range returns {
		returnedOrder := ReturnOrderAndCustomerToReturnedOrder(order)
		ordersProto[i] = &returnedOrder
	}

	return &pb.GetReturnedOrdersResponse{
		Orders: ordersProto,
	}, nil
}
