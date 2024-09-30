package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func (c *CliService) initInterCmd(rootCli *cobra.Command) {
	var interCli = &cobra.Command{
		Use:   "inter",
		Short: "Start program in interactive mode",
		Long:  `Start program in interactive mode`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.updateTimeInServiceInCmd(cmd); err != nil {
				fmt.Println(err.Error())
			}

			in := bufio.NewReader(os.Stdin)
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

	rootCli.AddCommand(interCli)

	interCli.Flags().StringP("start", "s", c.getToday(), "PVZ start time in format DD.MM.YYYY")
}
