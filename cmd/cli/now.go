package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *CliService) now() error {
	//if c.srvc == nil {
	//	return errors.New("this command should be use only in interactive mode")
	//}

	fmt.Println(c.getToday())
	return nil
}

//func init() {
//	rootCli.AddCommand(nowCli)
//
//	nowCli.AddCommand(&cobra.Command{
//		Use:   "help",
//		Short: "Help about command",
//		Long:  `Help about command`,
//		Run: func(cmd *cobra.Command, args []string) {
//			if err := nowCli.Help(); err != nil {
//				fmt.Println(err.Error())
//			}
//		},
//	})
//}

func (c *CliService) initNowCmd(rootCli *cobra.Command) {
	var nowCli = &cobra.Command{
		Use:     "now",
		Short:   "Get time in program",
		Long:    `Get time in program`,
		Example: "now",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.now(); err != nil {
				fmt.Println(err)
			}
		},
	}
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
