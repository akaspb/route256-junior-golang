package test_cli

import (
	"bytes"
	"context"
	"testing"
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers/storage"
)

func TestReturnsCmd(t *testing.T) {
	const correctOutputLen = 4

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
			Time: prevTime,
		},
	}

	order1 := baseOrder
	order1.ID = 1
	order1.Status.Value = models.StatusReturn
	order1.Expiry = prevTime

	order2 := baseOrder
	order2.ID = 2
	order2.Status.Value = models.StatusReturn
	order2.Expiry = nowTime

	order3 := baseOrder
	order3.ID = 3
	order3.Status.Value = models.StatusReturn
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

	var buffer bytes.Buffer
	err = cli.ReturnsHandler(ctx, &buffer, service, 0, 10)

	output := buffer.String()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	outputNewlinesCount := helpers.CountNewlines(output)
	if outputNewlinesCount != correctOutputLen {
		t.Errorf("wrong lines count (correct is %v):\n%v", outputNewlinesCount, output)
	}
}
