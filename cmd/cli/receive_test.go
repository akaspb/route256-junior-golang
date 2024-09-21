package cli

import (
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/test/cli/helpers"
	"testing"
	"time"
)

func TestReceive(t *testing.T) {
	tests := []struct {
		testName    string
		orderExpiry time.Time
		packName    string
		isErr       bool
	}{
		{
			testName:    "expired_already",
			orderExpiry: todayTime.Add(-day),
			packName:    "wrap",
			isErr:       true,
		},
		{
			testName:    "no_error",
			orderExpiry: todayTime.Add(day),
			packName:    "wrap",
			isErr:       false,
		},
		{
			testName:    "no_such_pack",
			orderExpiry: todayTime.Add(day),
			packName:    "no_such_pack",
			isErr:       true,
		},
	}

	for _, tc := range tests {
		jsonFile := fmt.Sprintf("/TestReceive_%s_storage.json", tc.testName)
		if err := helpers.DeleteFile(jsonFile); err != nil {
			t.Fatalf("unexpected error before test: %v", err)
		}
		//res, err := helpers.ExecuteCliCommand(t, now, tc.args...)

		//if res != "" {
		//	t.Errorf("unexpected error: %v", res)
		//}

		cliService := NewCliService(jsonFile)
		if err := cliService.getService(todayTime); err != nil {
			t.Fatalf("unexpected error before test: %v", err)
		}

		err := cliService.receive(
			0,
			1,
			1,
			tc.orderExpiry,
			tc.packName,
			0,
		)

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
