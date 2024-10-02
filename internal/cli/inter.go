package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
)

func getInterCmd(service *srvc.Service, rootCli *cobra.Command) *cobra.Command {
	var interCli = &cobra.Command{
		Use:   "inter",
		Short: "Start program in interactive mode",
		Long:  `Start program in interactive mode`,
		Run: func(cmd *cobra.Command, args []string) {
			if startTime, err := getStartTimeInCmd(cmd); err != nil {
				if !errors.Is(err, ErrorNoStartTimeInCMD) {
					fmt.Println(err.Error())
					return
				}
			} else {
				service.SetStartTime(startTime)
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

	interCli.Flags().StringP("start", "s", getToday(service), "PVZ start time in format DD.MM.YYYY")

	return interCli
}
