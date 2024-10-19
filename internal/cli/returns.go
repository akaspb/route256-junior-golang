package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/internal/pvz-service/v1"
)

func getReturnsCmd(client pvz_service.PvzServiceClient) *cobra.Command {
	var returnsCli = &cobra.Command{
		Use:     "returns",
		Short:   "Get all orders, which must be given to courier/couriers for return from PVZ",
		Long:    `Get all orders, which must be given to courier/couriers for return from PVZ`,
		Example: "returns -o=<offset> -l=<limit>",
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flags().Changed("offset") {
				fmt.Println("offset flag is not defined, check 'returns --help'")
				return
			}
			offset, err := cmd.Flags().GetInt("offset")
			if err != nil {
				fmt.Println(err.Error())
			}

			if !cmd.Flags().Changed("limit") {
				fmt.Println("limit flag is not defined, check 'returns --help'")
				return
			}
			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			request := &pvz_service.GetReturnedOrdersRequest{
				Offset: uint32(offset),
				Limit:  uint32(limit),
			}

			response, err := client.GetReturnedOrders(cmd.Context(), request)
			if err != nil {
				handleResponseError(err)
				return
			}

			if len(response.GetOrders()) == 0 {
				fmt.Println("No orders to show with such offset and limit")
				return
			}

			fmt.Printf("%8s|%11s\n", "Order ID", "Customer ID")
			for _, order := range response.Orders {
				fmt.Printf("%8v|%11v\n", order.GetOrderId(), order.GetCustomerId())
			}
		},
	}

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
	returnsCli.Flags().IntP("offset", "o", 0, "offset of results, starts from 0")

	return returnsCli
}
