package server

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
)

var (
	EventReceiveOrder = event_logger.EventType("order-created")
	EventGiveOrders   = event_logger.EventType("order-canceled")
	EventReturnOrder  = event_logger.EventType("order-canceled")
)

type Implementation struct {
	service      *service.Service
	eventLogger  event_logger.EventLogger
	eventFactory event_logger.EventFactory

	pb.UnimplementedPvzServiceServer
}

func NewImplementation(
	service *service.Service,
	eventLogger event_logger.EventLogger,
	eventFactory event_logger.EventFactory,
) *Implementation {
	return &Implementation{
		service:      service,
		eventLogger:  eventLogger,
		eventFactory: eventFactory,
	}
}
