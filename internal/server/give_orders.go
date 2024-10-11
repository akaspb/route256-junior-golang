package server

import (
	"context"

	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GiveOrders(ctx context.Context, req *pb.GiveOrdersRequest) (*pb.GiveOrdersResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orders, err := s.service.GiveOrderToCustomer(
		ctx,
		int64SlcToIDTypeSlc(req.GetOrderIds()),
		int64ToIDType(req.GetCustomerId()),
	)
	if err != nil {
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
