package models

import (
	"time"
)

type IDType int64
type WeightType float32
type CostType float32

type Order struct {
	ID         IDType     `json:"id"`
	CustomerID IDType     `json:"customerID"`
	Expiry     time.Time  `json:"expiry"`
	Weight     WeightType `json:"weight"`
	Cost       CostType   `json:"cost"`
	Packaging  *Packaging `json:"packaging"`
	Status     `json:"status"`
}

type StatusVal string

const (
	StatusToStorage  = StatusVal("to_storage")
	StatusToCustomer = StatusVal("to_customer")
	StatusReturn     = StatusVal("return")
)

type Status struct {
	Value StatusVal `json:"value"`
	Time  time.Time `json:"time"`
}

type Packaging struct {
	Name           string     `json:"name"`
	Cost           CostType   `json:"cost"`
	MaxOrderWeight WeightType `json:"maxOrderWeight"`
}
