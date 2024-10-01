package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func now() error {
	//if c.srvc == nil {
	//	return errors.New("this command should be use only in interactive mode")
	//}

	fmt.Println(getToday())
	return nil
}

func getNowCmd() *cobra.Command {
	var nowCli = &cobra.Command{
		Use:     "now",
		Short:   "Get time in program",
		Long:    `Get time in program`,
		Example: "now",
		Run: func(cmd *cobra.Command, args []string) {
			if err := now(); err != nil {
				fmt.Println(err)
			}
		},
	}

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

	return nowCli
}
