package server

import (
	"context"

	desc "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GiveOrders(ctx context.Context, req *desc.GiveOrdersRequest) (*desc.GiveOrdersResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orders, err := i.service.GiveOrderToCustomer(
		ctx,
		int64SlcToIDTypeSlc(req.OrderIds),
		int64ToIDType(req.CustomerId),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.GiveOrdersResponse{
		Orders: OrderIDWithMsgSlcToOrderToGiveInfoSlc(orders),
	}, nil
}
