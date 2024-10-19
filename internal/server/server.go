package server

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	pb "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
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
