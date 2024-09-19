package packaging

import (
	"errors"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

func NewPackaging(name string, cost models.CostType, maxOrderWeight models.WeightType) *models.Packaging {
	return &models.Packaging{
		Name:           name,
		Cost:           cost,
		MaxOrderWeight: maxOrderWeight,
	}
}

const (
	AnyWeight = -1
)

var Package = NewPackaging("package", 5, 10)
var Box = NewPackaging("box", 20, 30)
var Wrap = NewPackaging("wrap", 1, AnyWeight)

var Packs = map[string]*models.Packaging{
	Package.Name: Package,
	Box.Name:     Box,
	Wrap.Name:    Wrap,
}

func GetPackagingByName(packagingName string) (*models.Packaging, error) {
	if pack, ok := Packs[packagingName]; ok {
		return pack, nil
	}

	return nil, errors.New("no packaging with such name")
}
