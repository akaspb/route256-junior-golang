package postgres

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"time"
)

type Order struct {
	ID         models.IDType     `db:"id"`
	CustomerID models.IDType     `db:"customer_id"`
	Expiry     time.Time         `db:"expiry"`
	Weight     models.WeightType `db:"weight"`
	Cost       models.CostType   `db:"cost"`
	StatusID   models.IDType     `db:"status_id"`
}

type Status struct {
	ID    models.IDType    `db:"id"`
	Value models.StatusVal `db:"value"`
	Time  time.Time        `db:"time"`
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
