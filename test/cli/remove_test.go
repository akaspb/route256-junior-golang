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

func TestRemoveCmd(t *testing.T) {
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
			ChangedAt: prevTime,
		},
	}

	order1 := baseOrder
	order1.ID = 1
	order1.Status.Value = models.StatusToStorage
	order1.Expiry = prevTime

	order2 := baseOrder
	order2.ID = 2
	order2.Status.Value = models.StatusReturn
	order2.Expiry = futureTime

	order3 := baseOrder
	order3.ID = 3
	order3.Status.Value = models.StatusToCustomer
	order3.Expiry = futureTime

	someUnknownOrder := baseOrder
	someUnknownOrder.ID = 4
	someUnknownOrder.Status.Value = models.StatusToStorage
	someUnknownOrder.Expiry = prevTime

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
		order  models.Order
		output string
		err    error
	}{
		{
			order: order1,
			output: helpers.KeepСhars(`
				success: order can be given to courier for return
			`),
		},
		{
			order: order2,
			output: helpers.KeepСhars(`
				success: order can be given to courier for return
			`),
		},
		{
			order: order3,
			err:   srvc.ErrorOrderWasTakenByCustomer,
		},
		{
			order: someUnknownOrder,
			err:   srvc.ErrorOrderWasNotFounded,
		},
	}

	for _, tc := range tests {
		var buffer bytes.Buffer
		err = cli.RemoveHandler(ctx, &buffer, service, tc.order.ID)

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
