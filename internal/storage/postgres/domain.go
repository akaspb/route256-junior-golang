package postgres

import (
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

type Order struct {
	ID         models.IDType     `db:"id"`
	CustomerID models.IDType     `db:"customer_id"`
	Expiry     time.Time         `db:"expiry"`
	Weight     models.WeightType `db:"weight"`
	Cost       models.CostType   `db:"cost"`
}

type Status struct {
	OrderID   models.IDType    `db:"order_id"`
	Value     models.StatusVal `db:"value"`
	ChangedAt time.Time        `db:"changed_at"`
}

type Packaging struct {
	ID             models.IDType     `db:"id"`
	Name           string            `db:"name"`
	Cost           models.CostType   `db:"cost"`
	MaxOrderWeight models.WeightType `db:"max_order_weight"`
}

type Pack struct {
	OrderID        models.IDType     `db:"order_id"`
	Name           string            `db:"name"`
	Cost           models.CostType   `db:"cost"`
	MaxOrderWeight models.WeightType `db:"max_order_weight"`
}
