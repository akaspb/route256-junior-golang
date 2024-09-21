package cli

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/test/cli/helpers"
	"testing"
)

func TestGive(t *testing.T) {
	jsonFileCopy := "/TestGive_storage.json"
	if err := helpers.DeleteFile(jsonFileCopy); err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}
	if err := helpers.CopyFile("./test_storages/give_test_storage.json", jsonFileCopy); err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	cliService := NewCliService(jsonFileCopy)
	if err := cliService.getService(todayTime); err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	tests := []struct {
		testName   string
		orderId    []models.IDType
		customerId models.IDType
		isErr      bool
	}{
		{
			testName:   "another_order_customer",
			orderId:    []models.IDType{0},
			customerId: 1,
			isErr:      true,
		},
		{
			testName:   "expired",
			orderId:    []models.IDType{0},
			customerId: 0,
			isErr:      false,
		},
		{
			testName:   "ok",
			orderId:    []models.IDType{1},
			customerId: 0,
			isErr:      false,
		},
		{
			testName:   "was_taken_by_customer",
			orderId:    []models.IDType{2},
			customerId: 0,
			isErr:      false,
		},
		{
			testName:   "was_returned_by_customer",
			orderId:    []models.IDType{3},
			customerId: 0,
			isErr:      false,
		},
	}

	for _, tc := range tests {
		err := cliService.give(tc.customerId, tc.orderId)
		if tc.isErr {
			if err == nil {
				t.Errorf("error was expected, but give result: %v", err)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}
