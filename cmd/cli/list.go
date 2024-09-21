package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

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

		n, err := cmd.Flags().GetUint("n")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		customerIDint64, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		customerID := models.IDType(customerIDint64)

		orders, err := cliService.srvc.GetCustomerOrders(customerID, n)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		if len(orders) == 0 {
			fmt.Println("No orders")
			return
		}

		maxPackageNameLen := len("Pack")
		for _, order := range orders {
			maxPackageNameLen = max(maxPackageNameLen, len(order.Package))
		}

		tableTop := fmt.Sprintf("%8s|    Expiry|Expired|%"+strconv.Itoa(maxPackageNameLen)+"s|Cost", "ID", "Pack")
		fmt.Println(tableTop)

		for _, order := range orders {
			expired := "NO"
			if order.Expired {
				expired = "YES"
			}

			tableRow := fmt.Sprintf(
				"%8v|%-10s|%7s|%"+strconv.Itoa(maxPackageNameLen)+"s|%v",
				order.ID, order.Expiry.Format("02.01.2006"), expired, order.Package, order.Cost,
			)
			fmt.Println(tableRow)
		}
	},
}

func init() {
	RootCli.AddCommand(listCli)

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
}
