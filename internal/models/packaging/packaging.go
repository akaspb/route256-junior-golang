package packaging

import "gitlab.ozon.dev/siralexpeter/Homework/internal/models"

type Packaging struct {
	name           string
	cost           models.CostType
	maxOrderWeight models.WeightType
}

func NewPackaging(name string, cost models.CostType, maxOrderWeight models.WeightType) Packaging {
	return Packaging{
		name:           name,
		cost:           cost,
		maxOrderWeight: maxOrderWeight,
	}
}

func (p Packaging) GetName() string {
	return p.name
}

func (p Packaging) GetCost() models.CostType {
	return p.cost
}

func (p Packaging) GetMaxOrderWeight() models.WeightType {
	return p.maxOrderWeight
}
