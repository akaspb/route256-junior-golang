package server

import (
	"context"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetReturnedOrders(ctx context.Context, req *pb.GetReturnedOrdersRequest) (*pb.GetReturnedOrdersResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	returnsChan, err := s.service.GetReturnsList(ctx, int(req.GetOffset()), int(req.GetLimit()))

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ordersProto := make([]*pb.ReturnedOrder, 0, len(returnsChan))
	for order := range returnsChan {
		returnedOrder := ReturnOrderAndCustomerToReturnedOrder(order)
		ordersProto = append(ordersProto, &returnedOrder)
	}

	return &pb.GetReturnedOrdersResponse{
		Orders: ordersProto,
	}, nil
}
