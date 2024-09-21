package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var nowCli = &cobra.Command{
	Use:     "now",
	Short:   "Get time in program",
	Long:    `Get time in program`,
	Example: "now",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cliService.now(); err != nil {
			fmt.Println(err)
		}
	},
}

func (s *CliService) now() error {
	if s.srvc == nil {
		return errors.New("this command should be use only in interactive mode")
	}

	fmt.Println(s.srvc.GetCurrentTime().Format("02.01.2006"))
	return nil
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
