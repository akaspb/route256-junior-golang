package server

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	desc "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
)

type Implementation struct {
	service service.Service

	desc.UnimplementedPvzServiceServer
}

func NewImplementation(service service.Service) *Implementation {
	return &Implementation{service: service}
}
