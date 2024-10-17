package server

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
)

type Implementation struct {
	service     *service.Service
	eventLogger event_logger.Facade

	pb.UnimplementedPvzServiceServer
}

func NewImplementation(service *service.Service, eventLogger event_logger.Facade) *Implementation {
	return &Implementation{
		service:     service,
		eventLogger: eventLogger,
	}
}
