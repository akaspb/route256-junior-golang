package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var nowCli = &cobra.Command{
	Use:     "now",
	Short:   "Get time in program",
	Long:    `Get time in program`,
	Example: "now",
	Run: func(cmd *cobra.Command, args []string) {
		if cliService.srvc == nil {
			fmt.Println("this command should be use only in interactive mode")
			return
		}

		fmt.Println(cliService.srvc.GetCurrentTime().Format("02.01.2006"))
	},
}

func init() {
	rootCli.AddCommand(nowCli)

	nowCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := nowCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})
}
