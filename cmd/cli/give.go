package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"strconv"
)

var giveCli = &cobra.Command{
	Use:     "give",
	Short:   "Give orders by their ids from PVZ to customer",
	Long:    `Give orders by their ids from PVZ to customer`,
	Example: "give <userID> <orderID_1> ... <orderID_N>",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cliService.getServiceInCommand(cmd); err != nil {
			fmt.Println(err.Error())
			return
		}

		if len(args) < 2 {
			fmt.Println("IDs count is less then 2, check 'give --help'")
			return
		}

		ids := make([]models.IDType, len(args))
		for i, idStr := range args {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			ids[i] = models.IDType(id)
		}

		orders, err := cliService.srvc.GiveOrderToCustomer(ids[1:], ids[0])
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}

		if len(orders) == 0 {
			fmt.Println("No orders")
			return
		}

		maxMsgLen := 0
		maxPackageNameLen := len("Packaging")
		for _, order := range orders {
			maxMsgLen = max(maxMsgLen, len(order.Msg))
			maxPackageNameLen = max(maxPackageNameLen, len(order.Package))
		}

		tableTop := fmt.Sprintf(
			"%8s|Give|%-"+strconv.Itoa(maxMsgLen)+"s|%"+strconv.Itoa(maxPackageNameLen)+"s|Cost",
			"ID",
			"Message",
			"Packaging",
		)
		fmt.Println(tableTop)
		for _, order := range orders {
			give := "NO"
			if order.Ok {
				give = "YES"
			}
			tableRow := fmt.Sprintf(
				"%8v|%4s|%-"+strconv.Itoa(maxMsgLen)+"s|%-"+strconv.Itoa(maxPackageNameLen)+"s|%v",
				order.ID, give, order.Msg, order.Package, order.Cost)
			fmt.Println(tableRow)
		}
	},
}

func init() {
	rootCli.AddCommand(giveCli)

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
}
