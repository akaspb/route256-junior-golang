package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
)

//func GiveHandler(ctx context.Context, buffer *bytes.Buffer, service *srvc.Service, customerID models.IDType, orderIDs []models.IDType) error {
//	orders, err := service.GiveOrderToCustomer(ctx, orderIDs, customerID)
//	if err != nil {
//		return fmt.Errorf("error: %w", err)
//	}
//
//	if len(orders) == 0 {
//		fmt.Fprintln(buffer, "No orders")
//		return nil
//	}
//
//	maxMsgLen := 0
//	maxPackageNameLen := len("Pack")
//	for _, order := range orders {
//		maxMsgLen = max(maxMsgLen, len(order.Msg))
//		maxPackageNameLen = max(maxPackageNameLen, len(order.Package))
//	}
//
//	tableTop := fmt.Sprintf(
//		"%8s|Give|%-"+strconv.Itoa(maxMsgLen)+"s|%"+strconv.Itoa(maxPackageNameLen)+"s|Cost",
//		"ID",
//		"Message",
//		"Pack",
//	)
//	fmt.Fprintln(buffer, tableTop)
//	for _, order := range orders {
//		give := "NO"
//		if order.Ok {
//			give = "YES"
//		}
//		tableRow := fmt.Sprintf(
//			"%8v|%4s|%-"+strconv.Itoa(maxMsgLen)+"s|%-"+strconv.Itoa(maxPackageNameLen)+"s|%v",
//			order.ID, give, order.Msg, order.Package, order.Cost)
//		fmt.Fprintln(buffer, tableRow)
//	}
//	return nil
//}

func getGiveCmd(client pvz_service.PvzServiceClient) *cobra.Command {
	var giveCli = &cobra.Command{
		Use:     "give",
		Short:   "Give orders by their ids from PVZ to customer",
		Long:    `Give orders by their ids from PVZ to customer`,
		Example: "give <userID> <orderID_1> ... <orderID_N>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				fmt.Println("IDs count is less then 2, check 'GiveHandler --help'")
				return
			}

			ids := make([]int64, len(args))
			for i, idStr := range args {
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				ids[i] = id
			}

			customerID := ids[0]
			orderIDs := ids[1:]

			request := &pvz_service.GiveOrdersRequest{
				CustomerId: customerID,
				OrderIds:   orderIDs,
			}

			response, err := client.GiveOrders(cmd.Context(), request)
			if err != nil {
				handleResponseError(err)
				return
			}

			maxMsgLen := 0
			maxPackageNameLen := len("Pack")
			for _, order := range response.GetOrders() {
				maxMsgLen = max(maxMsgLen, len(order.GetMessage()))
				maxPackageNameLen = max(maxPackageNameLen, len(order.GetPacking()))
			}

			if len(response.GetOrders()) == 0 {
				fmt.Println("No orders for this customer")
			}

			fmt.Printf(
				"%8s|Give|%-"+strconv.Itoa(maxMsgLen)+"s|%"+strconv.Itoa(maxPackageNameLen)+"s|Cost\n",
				"ID",
				"Message",
				"Pack",
			)

			for _, order := range response.GetOrders() {
				give := "NO"
				if order.GetOkToGive() {
					give = "YES"
				}
				fmt.Printf(
					"%8v|%4s|%-"+strconv.Itoa(maxMsgLen)+"s|%-"+strconv.Itoa(maxPackageNameLen)+"s|%v\n",
					order.GetOrderId(), give, order.GetMessage(), order.GetPacking(), order.Cost)
			}
		},
	}

	giveCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := giveCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})

	return giveCli
}
