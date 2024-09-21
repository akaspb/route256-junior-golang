package cli

import (
	"testing"
)

func TestNowCmdCommand(t *testing.T) {
	cliServiceFullInited := NewCliService("/TestNowCmdCommand_storage.json")
	if err := cliServiceFullInited.getService(todayTime); err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	tests := []struct {
		cliService *CliService
		isErr      bool
	}{
		{
			cliService: NewCliService(""),
			isErr:      true,
		},
		{
			cliService: cliServiceFullInited,
			isErr:      false,
		},
	}

	for _, tc := range tests {
		//res, err := helpers.ExecuteCliCommand(t, now, tc.args...)

		//if res != "" {
		//	t.Errorf("unexpected error: %v", res)
		//}

		err := tc.cliService.now()

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
