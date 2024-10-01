package test_cli

import (
	"bytes"
	"context"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers/storage"
	"testing"
	"time"
)

func TestListCmd(t *testing.T) {
	ctx := context.Background()
	testStorage := storage.NewStorage()

	prevTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	nowTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	futureTime := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	orderOk := models.Order{
		ID:         1,
		CustomerID: 1,
		Expiry:     futureTime,
		Weight:     1,
		Cost:       1,
		Pack:       nil,
		Status: models.Status{
			Value: models.StatusToStorage,
			Time:  prevTime,
		},
	}

	orderExpired := models.Order{
		ID:         2,
		CustomerID: 2,
		Expiry:     prevTime,
		Weight:     1,
		Cost:       1,
		Pack:       nil,
		Status: models.Status{
			Value: models.StatusToStorage,
			Time:  prevTime,
		},
	}

	err := testStorage.FillWithOrders(ctx, orderOk, orderExpired)
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
		output     string
	}{
		{
			customerID: 1,
			output: helpers.KeepСhars(`
              ID|    Expiry|Expired|Pack|Cost
               1|01.01.2022|     NO|    |1
			`),
		},
		{
			customerID: 2,
			output: helpers.KeepСhars(`
              ID|    Expiry|Expired|Pack|Cost
               2|01.01.2020|    YES|    |1
			`),
		},
	}

	for _, tc := range tests {
		var buffer bytes.Buffer
		err = cli.ListHandler(ctx, &buffer, service, tc.customerID, 1)

		output := buffer.String()
		output = helpers.KeepСhars(output)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		if output != tc.output {
			t.Errorf("wrong output: %v", output)
		}
	}
}
