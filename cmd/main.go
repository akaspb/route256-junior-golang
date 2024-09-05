package main

import (
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"time"
)

const (
	jsonPath = "internal/storage/storage.json"
)

func main() {
	orderStorage, err := storage.InitJsonStorage(jsonPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	currTime := time.Now()
	orderService := service.NewService(orderStorage, currTime, currTime)

	//err = orderService.AcceptOrderFromCourier(7, 1, currTime.Add(time.Second*2))
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//err = orderService.AcceptOrderFromCourier(8, 1, currTime.Add(time.Second*2))
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//err = orderService.AcceptOrderFromCourier(9, 1, currTime.Add(time.Second*2))
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	fmt.Println(orderStorage.OrderStorage)

	//orderIDs, err := orderService.GetCustomerOrders(1, 0)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(orderIDs)

	err = orderService.ReturnOrder(3)
	if err != nil {
		fmt.Println(err.Error())
	}

	orderIDs, err := orderService.GetReturnsList()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(orderIDs)
	}

	//time.Sleep(time.Second * 1)
	//
	//ordersToGive, err := orderService.GiveOrderToCustomer([]models.IDType{6}, 2)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(ordersToGive)

	//time.Sleep(time.Second * 3)

	//err = orderService.ReturnOrder(3)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	//orderIDs, err := orderService.GetCustomerOrders(1, 0)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(orderIDs)
	//
	//err = orderService.ReturnOrder(3)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//orderIDs, err = orderService.GetCustomerOrders(1, 0)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(orderIDs)

	//ordersToGive, err := orderService.GiveOrderToCustomer([]models.IDType{1}, 1)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(ordersToGive)
}
