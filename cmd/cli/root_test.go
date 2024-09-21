package cli_test

import (
	"gitlab.ozon.dev/siralexpeter/Homework/cmd/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/test/cli/helpers"
	"testing"
)

func TestRootCmdCommand(t *testing.T) {
	root := cli.RootCli

	tests := []struct {
		args  []string
		isErr bool
	}{
		{
			args:  nil,
			isErr: false,
		},
		//{
		//	args: nil,
		//	err:  errors.New("not ok"),
		//},
		//{
		//	args: []string{"-t"},
		//	err:  nil,
		//	out:  "ok",
		//},
		//{
		//	args: []string{"--toggle"},
		//	err:  nil,
		//	out:  "ok",
		//},
	}

	for _, tc := range tests {
		res, err := helpers.ExecuteCliCommand(t, root, tc.args...)

		if tc.isErr {
			if err == nil {
				t.Errorf("error was expected, but give result: %v", res)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}
