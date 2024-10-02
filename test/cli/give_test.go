package test_cli

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers/storage"
)

func TestGiveCmd(t *testing.T) {
	ctx := context.Background()
	testStorage := storage.NewStorage()

	prevTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	nowTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	futureTime := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	baseOrder := models.Order{
		ID:         0,
		CustomerID: 0,
		Expiry:     futureTime,
		Weight:     1,
		Cost:       1,
		Pack:       nil,
		Status: models.Status{
			Value: models.StatusToStorage,
			Time:  prevTime,
		},
	}

	orderUser1 := baseOrder
	orderUser1.ID = 1
	orderUser1.CustomerID = 1

	orderUser2 := baseOrder
	orderUser2.ID = 2
	orderUser2.CustomerID = 2

	err := testStorage.FillWithOrders(ctx, orderUser1, orderUser2)
	if err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	packService, err := packaging.NewPackaging()
	if err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	service := srvc.NewService(testStorage, packService, nowTime, time.Now())

	tests := []struct {
		customerID models.IDType
		orderIDs   []models.IDType
		output     string
		err        error
	}{
		{
			customerID: 1,
			orderIDs:   []models.IDType{1},
			output: helpers.KeepСhars(`
				ID|Give|Message               |Pack|Cost
				 1| YES|Give order to customer|    |1
			`),
		},
		{
			customerID: 1,
			orderIDs:   []models.IDType{1, 2},
			err:        srvc.ErrorCustomerID,
		},
	}

	for _, tc := range tests {
		var buffer bytes.Buffer
		err = cli.GiveHandler(ctx, &buffer, service, tc.customerID, tc.orderIDs)

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
