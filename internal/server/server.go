package server

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	desc "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
)

type Implementation struct {
	storage storage.Facade

	desc.UnimplementedPvzServiceServer
}
