package server

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func int64ToIDType(id int64) models.IDType {
	return models.IDType(id)
}

func idTypeToInt64(id models.IDType) int64 {
	return int64(id)
}

func costTypeToFloat(cost models.CostType) float32 {
	return float32(cost)
}

func floatToCostType(cost float32) models.CostType {
	return models.CostType(cost)
}

func weightTypeToFloat(weight models.WeightType) float32 {
	return float32(weight)
}

func floatToWeightType(weight float32) models.WeightType {
	return models.WeightType(weight)
}

func orderIDWithMsgToOrderToGiveInfo(order service.OrderIDWithMsg) pb.OrderToGiveInfo {
	return pb.OrderToGiveInfo{
		OrderId:  idTypeToInt64(order.ID),
		Cost:     costTypeToFloat(order.Cost),
		Packing:  wrapperspb.String(order.Package),
		Message:  order.Msg,
		OkToGive: order.Ok,
	}
}

func orderIDWithExpiryAndStatusToCustomerOrderInfo(order service.OrderIDWithExpiryAndStatus) pb.CustomerOrderInfo {
	return pb.CustomerOrderInfo{
		OrderId: idTypeToInt64(order.ID),
		Weight:  weightTypeToFloat(order.Weight),
		Cost:    costTypeToFloat(order.Cost),
		Expiry:  timestamppb.New(order.Expiry),
		Expired: order.Expired,
		Packing: order.Package,
	}
}

func ReturnOrderAndCustomerToReturnedOrder(order service.ReturnOrderAndCustomer) pb.ReturnedOrder {
	return pb.ReturnedOrder{
		OrderId:    idTypeToInt64(order.OrderID),
		CustomerId: idTypeToInt64(order.CustomerID),
	}
}
