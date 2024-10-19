package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/internal/pvz-service/v1"
)

func getReturnCmd(client pvz_service.PvzServiceClient) *cobra.Command {
	var returnCli = &cobra.Command{
		Use:     "return",
		Short:   "Get order from customer to return",
		Long:    `Get order from customer to return`,
		Example: "return -o=<orderID> -c=<customerID>",
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flags().Changed("order") {
				fmt.Println("order flag is not defined, check 'return --help'")
				return
			}
			orderID, err := cmd.Flags().GetInt64("order")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !cmd.Flags().Changed("customer") {
				fmt.Println("customer flag is not defined, check 'return --help'")
				return
			}
			customerID, err := cmd.Flags().GetInt64("customer")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			request := &pvz_service.ReturnOrderRequest{
				CustomerId: customerID,
				OrderId:    orderID,
			}

			_, err = client.ReturnOrder(cmd.Context(), request)
			if err != nil {
				handleResponseError(err)
				return
			}

			fmt.Println("success")
		},
	}

	returnCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := returnCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})

	returnCli.Flags().Int64P("order", "o", 0, "unique order ID")
	returnCli.Flags().Int64P("customer", "c", 0, "unique customer ID")

	return returnCli
}
