package cli

import "gitlab.ozon.dev/go/classroom-15/students/workshop-1/internal/models"

func AcceptOrderFromCourier(orderID, customerID models.IDType, orderExpiry models.Time) error {

	return nil
}

func ReturnOrderToCourier(orderID models.IDType) error {

}

func GiveOrderToCustomer(orderIDs []models.IDType) error {

}

func GetCustomerOrders(customerID models.IDType) ([]models.IDType, error) {

}

func ReturnOrderFromCustomer(customerID, orderID models.IDType) error {

}

func GetReturnsList() ([]models.IDType, error) {

}
