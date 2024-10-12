package server

import (
	"context"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) RemoveOrder(ctx context.Context, req *pb.RemoveOrderRequest) (*pb.RemoveOrderResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.ReturnOrder(ctx, int64ToIDType(req.GetOrderId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveOrderResponse{}, nil
}
