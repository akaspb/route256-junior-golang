package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

var removeCli = &cobra.Command{
	Use:     "remove",
	Short:   "Return order from PVZ to courier",
	Long:    `Return order from PVZ to courier`,
	Example: "remove <orderID>",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cliService.getServiceInCommand(cmd); err != nil {
			fmt.Println(err.Error())
			return
		}

		if len(args) < 1 {
			fmt.Println("orderID is not defined, check 'remove --help'")
			return
		}

		orderIDint64, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		orderID := models.IDType(orderIDint64)

		if err := cliService.srvc.ReturnOrder(orderID); err != nil {
			fmt.Printf("error: %v\n", err)
		} else {
			fmt.Println("success: order can be given to courier for return")
		}
	},
}

func init() {
	rootCli.AddCommand(removeCli)

	removeCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := removeCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})
}
