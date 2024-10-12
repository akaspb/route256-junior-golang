package server

import (
	"context"
	"errors"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) ChangeThreadCount(ctx context.Context, req *pb.ChangeThreadCountRequest) (*pb.ChangeThreadCountResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.ChangeWorkerCount(
		int(req.GetThreadCount()),
	)
	if err != nil {
		var argumentErr *service.ArgumentError
		if errors.As(err, &argumentErr) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ChangeThreadCountResponse{}, nil
}
