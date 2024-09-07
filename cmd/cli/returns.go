package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var returnsCli = &cobra.Command{
	Use:     "returns",
	Short:   "Get all orders, which must be given to courier/couriers for return from PVZ",
	Long:    `Get all orders, which must be given to courier/couriers for return from PVZ`,
	Example: "returns -o=<offset> -l=<limit>",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cliService.getServiceInCommand(cmd); err != nil {
			fmt.Println(err.Error())
			return
		}

		var offset, limit int
		var err error

		if cmd.Flags().Changed("offset") {
			offset, err = cmd.Flags().GetInt("offset")
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("offset flag is not defined")
			return
		}

		if cmd.Flags().Changed("limit") {
			limit, err = cmd.Flags().GetInt("limit")
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("limit flag is not defined")
			return
		}

		if orders, err := cliService.srvc.GetReturnsList(offset, limit); err != nil {
			fmt.Println("error:", err.Error())
		} else {
			if len(orders) == 0 {
				fmt.Println("No orders")
			} else {
				tableTop := fmt.Sprintf("%8s|%11s", "Order ID", "Customer ID")
				fmt.Println(tableTop)
				for _, order := range orders {
					tableRow := fmt.Sprintf("%8v|%11v", order.OrderID, order.CustomerID)
					fmt.Println(tableRow)
				}
			}
		}
	},
}

func init() {
	rootCli.AddCommand(returnsCli)

	returnsCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := returnsCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})

	returnsCli.Flags().IntP("limit", "l", 0, "limit of results")
	returnsCli.Flags().IntP("offset", "o", 0, "offset of results")
}
