package server

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	desc "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func int64ToIDType(id int64) models.IDType {
	return models.IDType(id)
}

func idTypeToInt64(id models.IDType) int64 {
	return int64(id)
}

func int64SlcToIDTypeSlc(ids []int64) []models.IDType {
	res := make([]models.IDType, len(ids))
	for j, id := range ids {
		res[j] = int64ToIDType(id)
	}

	return res
}

func costTypeToFloat(cost models.CostType) float32 {
	return float32(cost)
}

func OrderIDWithMsgToOrderToGiveInfo(order service.OrderIDWithMsg) desc.OrderToGiveInfo {
	return desc.OrderToGiveInfo{
		OrderId:  idTypeToInt64(order.ID),
		Cost:     costTypeToFloat(order.Cost),
		Packing:  wrapperspb.String(order.Package),
		Message:  order.Msg,
		OkToGive: order.Ok,
	}
}

func OrderIDWithMsgSlcToOrderToGiveInfoSlc(orders []service.OrderIDWithMsg) []*desc.OrderToGiveInfo {
	res := make([]*desc.OrderToGiveInfo, len(orders))
	for i, order := range orders {
		orderToGiveInfo := OrderIDWithMsgToOrderToGiveInfo(order)
		res[i] = &orderToGiveInfo
	}

	return res
}
