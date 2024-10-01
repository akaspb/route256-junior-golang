package test_cli

import (
	"bytes"
	"context"
	"errors"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers/storage"
	"testing"
	"time"
)

func TestReturnCmd(t *testing.T) {
	ctx := context.Background()
	testStorage := storage.NewStorage()

	prevTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	nowTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	futureTime := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	baseOrder := models.Order{
		ID:         0,
		CustomerID: 1,
		Weight:     1,
		Cost:       1,
		Status: models.Status{
			Time:  nowTime,
			Value: models.StatusToCustomer,
		},
	}

	order1 := baseOrder
	order1.ID = 1
	order1.Expiry = futureTime

	order2 := baseOrder
	order2.ID = 2
	order2.Status.Time = prevTime

	order3 := baseOrder
	order3.ID = 3
	order3.Expiry = futureTime

	err := testStorage.FillWithOrders(ctx, order1, order2, order3)
	if err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	packService, err := packaging.NewPackaging()
	if err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	service := srvc.NewService(testStorage, packService, nowTime, time.Now())

	tests := []struct {
		orderID    models.IDType
		customerID models.IDType
		output     string
		err        error
	}{
		{
			orderID:    1,
			customerID: 1,
			output: helpers.KeepСhars(`
				success: take order from customer to store it in PVZ
			`),
		},
		{
			orderID:    2,
			customerID: 1,
			err:        srvc.ErrorOrderExpiredAlready,
		},
		{
			orderID:    3,
			customerID: 2,
			err:        srvc.ErrorCustomerID,
		},
		{
			orderID:    4,
			customerID: 1,
			err:        srvc.ErrorOrderWasNotFounded,
		},
	}

	for _, tc := range tests {
		var buffer bytes.Buffer
		err = cli.ReturnHandler(ctx, &buffer, service, tc.customerID, tc.orderID)

		output := buffer.String()
		output = helpers.KeepСhars(output)

		if err != nil {
			if !errors.Is(err, tc.err) {
				t.Errorf("unexpected error: %v", err)
			}
		} else {
			if tc.err != nil {
				t.Errorf("error was expected: %v", tc.err)
				continue
			}
			if output != tc.output {
				t.Errorf("wrong output: %v", output)
			}
		}
	}
}
