package server

import (
	"context"
	"errors"
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) ReceiveOrder(ctx context.Context, req *pb.ReceiveOrderRequest) (*pb.ReceiveOrderResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	s.logMethodCall(
		"ReceiveOrder",
		EventReceiveOrderReq,
		fmt.Sprintf("receive order id==%v customer_id==%v", req.GetId(), req.GetCustomerId()),
	)

	var packPtr *models.Pack
	if req.GetPacking() != nil {
		pack, err := s.service.Packaging.GetPackagingByName(req.GetPacking().GetValue())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		packPtr = &pack
	}

	err := s.service.AcceptOrderFromCourier(ctx, service.AcceptOrderDTO{
		OrderID:     int64ToIDType(req.GetId()),
		OrderCost:   floatToCostType(req.GetCost()),
		OderWeight:  floatToWeightType(req.GetWeight()),
		CustomerID:  int64ToIDType(req.GetCustomerId()),
		Pack:        packPtr,
		OrderExpiry: req.GetExpiry().AsTime(),
	})
	if err != nil {
		s.logMethodCall(
			"ReceiveOrder",
			EventReceiveOrderRes,
			fmt.Sprintf("error: %v", err),
		)

		var argumentErr *service.ArgumentError
		if errors.As(err, &argumentErr) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logMethodCall(
		"ReceiveOrder",
		EventReceiveOrderRes,
		"success",
	)

	return &pb.ReceiveOrderResponse{}, nil
}
