package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"testing"
	"time"
)

const (
	day      = 24 * time.Hour
	year     = 365 * day
	longTime = year
)

func TestGetCliService(t *testing.T) {
	tests := []struct {
		testName              string
		currUserSpecifiedTime time.Time
		isErr                 bool
	}{
		{
			testName:              "time_now",
			currUserSpecifiedTime: todayTime,
			isErr:                 false,
		},
		{
			testName:              "time_past",
			currUserSpecifiedTime: todayTime.Add(-longTime),
			isErr:                 false,
		},
		{
			testName:              "time_future",
			currUserSpecifiedTime: todayTime.Add(longTime),
			isErr:                 false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			cliService := NewCliService(fmt.Sprintf("/TestCliService_%s_storage.json", tc.testName))
			err := cliService.updateTimeInService(tc.currUserSpecifiedTime)

			if tc.isErr {
				if err == nil {
					t.Error("error was expected")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetCliNowCmdCommand(t *testing.T) {
	tests := []struct {
		testName              string
		currUserSpecifiedTime time.Time
		args                  []string
		isErr                 bool
	}{
		{
			testName:              "no_start_flag",
			args:                  nil,
			currUserSpecifiedTime: todayTime,
			isErr:                 false,
		},
		{
			testName: "start_flag",
			args:     []string{"--start", todayStr},
			isErr:    false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			cmd := &cobra.Command{
				Use: "test",
				Run: func(cmd *cobra.Command, args []string) {},
			}
			cmd.Flags().StringP("start", "s", "", "")

			cmd.SetArgs(tc.args)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("unexpected error before test: %v", err)
			}

			cliService := NewCliService(fmt.Sprintf("/TestGetCliNowCmdCommand_%s_storage.json", tc.testName))
			err := cliService.updateTimeInServiceInCmd(cmd)

			if tc.isErr {
				if err == nil {
					t.Error("error was expected")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
