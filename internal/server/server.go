package server

import (
	pb "gitlab.ozon.dev/siralexpeter/Homework/internal/pvz-service/v1"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
)

type Implementation struct {
	service *service.Service

	pb.UnimplementedPvzServiceServer
}

func NewImplementation(
	service *service.Service,
) *Implementation {
	return &Implementation{
		service: service,
	}
}
