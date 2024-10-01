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
	"gitlab.ozon.dev/siralexpeter/Homework/test/storage"
	"testing"
	"time"
)

func TestReceiveCmd(t *testing.T) {
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
		Status: models.Status{
			Value: models.StatusToStorage,
			Time:  prevTime,
		},
	}

	order1 := baseOrder
	order1.ID = 1
	order1.CustomerID = 1

	order2 := baseOrder
	order2.ID = 2
	order2.CustomerID = 2

	err := testStorage.FillWithOrders(ctx, order1)
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
			order: order2,
			output: helpers.KeepСhars(`
				success: take order for storage in PVZ
			`),
		},
		{
			order: order1,
			err:   srvc.ErrorOrderWasAccepted,
		},
		{
			order: order1,
			err:   srvc.ErrorOrderWasAccepted,
		},
	}

	for _, tc := range tests {
		var buffer bytes.Buffer
		err = cli.ReceiveHandler(cli.ReceiveHandlerDTO{
			Ctx:         ctx,
			Buffer:      &buffer,
			Service:     service,
			OrderID:     tc.order.ID,
			OrderCost:   tc.order.Cost,
			OrderWeight: tc.order.Weight,
			OrderExpiry: tc.order.Expiry,
			PackName:    "",
			CustomerID:  tc.order.CustomerID,
		})

		output := buffer.String()
		output = helpers.KeepСhars(output)

		if err != nil {
			if !errors.Is(err, tc.err) {
				t.Errorf("unexpected error: %v", err)
			}
		} else {
			if output != tc.output {
				t.Fatalf("wrong output: %v", output)
			}
		}
	}
}
