package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/internal/pvz-service/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func getListCmd(client pvz_service.PvzServiceClient) *cobra.Command {
	var listCli = &cobra.Command{
		Use:     "list",
		Short:   "Get customer orders, which are contained in PVZ now",
		Long:    `Get customer orders, which are contained in PVZ now`,
		Example: "list <userID> [-n=<N last orders count>]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("customer ID is not given, check 'list --help'")
				return
			}

			customerID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var lastCount *wrapperspb.UInt32Value
			n, err := cmd.Flags().GetUint("n")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if n > 0 {
				lastCount = wrapperspb.UInt32(uint32(n))
			}

			request := &pvz_service.GetCustomerOrdersRequest{
				CustomerId: customerID,
				LastCount:  lastCount,
			}

			response, err := client.GetCustomerOrders(cmd.Context(), request)
			if err != nil {
				handleResponseError(err)
				return
			}

			if len(response.GetOrders()) == 0 {
				fmt.Println("No orders")
				return
			}

			maxPackageNameLen := len("Pack")
			for _, order := range response.GetOrders() {
				maxPackageNameLen = max(maxPackageNameLen, len(order.GetOrderInfo().GetPacking()))
			}

			fmt.Printf("%8s|    Expiry|Expired|Weight|%"+strconv.Itoa(maxPackageNameLen)+"s|Cost\n", "ID", "Pack")

			for _, order := range response.GetOrders() {
				expired := "NO"
				if order.GetExpired() {
					expired = "YES"
				}

				fmt.Printf(
					"%8v|%-10s|%7s|%6.2f|%"+strconv.Itoa(maxPackageNameLen)+"s|%v\n",
					order.GetOrderInfo().GetOrderId(),
					order.GetExpiry().AsTime().Format("02.01.2006"),
					expired,
					order.GetOrderInfo().GetWeight(),
					order.GetOrderInfo().GetPacking(),
					order.GetOrderInfo().GetCost(),
				)
			}
		},
	}

	listCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := listCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})

	listCli.Flags().UintP("n", "n", 0, "N last orders count")

	return listCli
}
