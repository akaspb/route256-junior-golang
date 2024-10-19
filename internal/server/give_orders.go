package server

import (
	"context"
	"errors"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	pb "gitlab.ozon.dev/siralexpeter/Homework/internal/pvz-service/v1"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GiveOrders(ctx context.Context, req *pb.GiveOrdersRequest) (*pb.GiveOrdersResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orderIDsProto := req.GetOrderIds()
	orderIDs := make([]models.IDType, len(orderIDsProto))
	for j, id := range orderIDsProto {
		orderIDs[j] = int64ToIDType(id)
	}

	orders, err := s.service.GiveOrderToCustomer(
		ctx,
		orderIDs,
		int64ToIDType(req.GetCustomerId()),
	)
	if err != nil {
		var argumentErr *service.ArgumentError
		if errors.As(err, &argumentErr) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	ordersProto := make([]*pb.OrderToGiveInfo, len(orders))
	for i, order := range orders {
		orderToGiveInfo := orderIDWithMsgToOrderToGiveInfo(order)
		ordersProto[i] = &orderToGiveInfo
	}

	return &pb.GiveOrdersResponse{
		Orders: ordersProto,
	}, nil
}
