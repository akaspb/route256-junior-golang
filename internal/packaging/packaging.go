package packaging

import (
	"errors"
	"fmt"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

func NewPack(name string, cost models.CostType, maxOrderWeight models.WeightType) models.Pack {
	return models.Pack{
		Name:           name,
		Cost:           cost,
		MaxOrderWeight: maxOrderWeight,
	}
}

const (
	AnyWeight = -1
)

type Packaging struct {
	Packs map[string]models.Pack
}

func NewPackaging(packs ...models.Pack) (*Packaging, error) {
	packsMap := make(map[string]models.Pack, len(packs))
	for _, pack := range packs {
		if _, ok := packsMap[pack.Name]; ok {
			return nil, fmt.Errorf("non unique name '%s' in packs param", pack.Name)
		}
		packsMap[pack.Name] = pack
	}

	return &Packaging{Packs: packsMap}, nil
}

//func DefaultPackaging() *Packaging {
//	packet := NewPack("packet", 5, 10)
//	box := NewPack("box", 20, 30)
//	wrap := NewPack("wrap", 1, AnyWeight)
//
//	packaging, _ := NewPackaging(packet, box, wrap)
//	return packaging
//}

func (p *Packaging) GetPackagingByName(packagingName string) (models.Pack, error) {
	if pack, ok := p.Packs[packagingName]; ok {
		return pack, nil
	}

	return models.Pack{}, errors.New("no packaging with such name")
}

func (p *Packaging) PackOrder(pack models.Pack, orderWeight models.WeightType) (cost models.CostType, err error) {
	if pack.MaxOrderWeight != AnyWeight && orderWeight >= pack.MaxOrderWeight {
		return 0, fmt.Errorf(
			"order weight==%v reached max packaging '%s' weight==%v",
			orderWeight,
			pack.Name,
			pack.MaxOrderWeight,
		)
	}

	return pack.Cost, nil
}
