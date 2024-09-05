package models

import (
	"time"
)

type IDType int64

//type Time time.Time

type Order struct {
	ID     IDType    `json:"id"`
	Expiry time.Time `json:"expiry"`
	Status `json:"status"`
}

type StatusVal string

const (
	StatusToStorage  = StatusVal("to_storage")
	StatusToCustomer = StatusVal("to_customer")
	StatusReturn     = StatusVal("return")
	//StatusReturnFromCustomer = StatusVal("return_from_customer")
)

type Status struct {
	Val        StatusVal `json:"val"`
	CustomerID IDType    `json:"customerID"`
	Time       time.Time `json:"time"`
}
