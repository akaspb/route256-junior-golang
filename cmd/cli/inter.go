package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var interCli = &cobra.Command{
	Use:   "inter",
	Short: "Start program in interactive mode",
	Long:  `Start program in interactive mode`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cliService.getServiceInCommand(cmd); err != nil {
			fmt.Println(err.Error())
		}

		for {
			fmt.Print("> ")

			res, _ := in.ReadString('\n')
			res = strings.TrimSpace(res)
			if res == "exit" {
				break
			}

			res = strings.Replace(res, "--help", "help", -1)
			res = strings.Replace(res, "-h", "help", -1)

			argsFromLine := strings.Split(res, " ")
			rootCli.SetArgs(argsFromLine)
			if err := rootCli.Execute(); err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println()
		}
	},
}

func init() {
	rootCli.AddCommand(interCli)

	interCli.Flags().StringP("start", "s", todayStr, "PVZ start time in format DD.MM.YYYY")
}
