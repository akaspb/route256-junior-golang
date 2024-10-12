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
	packsMap map[string]models.Pack
}

func NewPackaging() (*Packaging, error) {
	p := &Packaging{packsMap: make(map[string]models.Pack)}

	packet := NewPack("packet", 5, 10)
	box := NewPack("box", 20, 30)
	wrap := NewPack("wrap", 1, AnyWeight)

	if err := p.AddPacks(packet, box, wrap); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Packaging) AddPacks(packs ...models.Pack) error {
	for _, pack := range packs {
		if err := p.AddPack(pack); err != nil {
			return err
		}
	}

	return nil
}

func (p *Packaging) AddPack(pack models.Pack) error {
	if _, ok := p.packsMap[pack.Name]; ok {
		return fmt.Errorf("non unique name '%s' in packs param", pack.Name)
	}

	p.packsMap[pack.Name] = pack

	return nil
}

func (p *Packaging) GetAllPacks() []models.Pack {
	packs := make([]models.Pack, 0, len(p.packsMap))
	for _, pack := range p.packsMap {
		packs = append(packs, pack)
	}

	return packs
}

func (p *Packaging) GetPackagingByName(packagingName string) (models.Pack, error) {
	if pack, ok := p.packsMap[packagingName]; ok {
		return pack, nil
	}

	return models.Pack{}, errors.New("no packaging with such name")
}

func (p *Packaging) PackOrder(pack models.Pack, orderWeight models.WeightType) (cost models.CostType, err error) {
	if pack.MaxOrderWeight != AnyWeight && !(pack.MaxOrderWeight > 0) {
		return 0, errors.New("error pack: not pack.MaxOrderWeight > 0")
	}

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
