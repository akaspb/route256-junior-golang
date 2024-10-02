package cli

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
)

func NowHandler(buffer *bytes.Buffer, service *service.Service) error {
	//if c.srvc == nil {
	//	return errors.New("this command should be use only in interactive mode")
	//}

	fmt.Fprintln(buffer, getToday(service))
	return nil
}

func getNowCmd(service *service.Service) *cobra.Command {
	var nowCli = &cobra.Command{
		Use:     "now",
		Short:   "Get time in program",
		Long:    `Get time in program`,
		Example: "now",
		Run: func(cmd *cobra.Command, args []string) {
			var buffer bytes.Buffer
			if err := NowHandler(&buffer, service); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Print(buffer.String())
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
