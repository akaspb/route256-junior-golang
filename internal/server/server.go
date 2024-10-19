package server

import (
	"encoding/json"
	"log"

	"github.com/golang/protobuf/proto"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
)

var (
	EventReceiveOrderReq = event_logger.EventType("receive_order_req")
	EventReceiveOrderRes = event_logger.EventType("receive_order_res")

	EventGiveOrdersReq = event_logger.EventType("give_orders_req")
	EventGiveOrdersRes = event_logger.EventType("give_orders_res")

	EventReturnOrderReq = event_logger.EventType("return_order_req")
	EventReturnOrderRes = event_logger.EventType("return_order_res")
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

func (s *Implementation) logMethodCallProtoMessage(method string, eventType event_logger.EventType, protoMsg proto.Message) {
	bytes, err := json.Marshal(protoMsg)
	if err != nil {
		handleLoggerError(method, err)
	}

	s.logMethodCall(method, eventType, string(bytes))
}

func (s *Implementation) logMethodCall(method string, eventType event_logger.EventType, details string) {
	event, err := s.eventFactory.Create(
		eventType,
		details,
	)
	if err != nil {
		handleLoggerError(method, err)
	}

	err = s.eventLogger.Send(event)
	if err != nil {
		handleLoggerError(method, err)
	}
}

func handleLoggerError(method string, err error) {
	log.Printf("[kafka producer] method: %s; error: %v", method, err)
}
