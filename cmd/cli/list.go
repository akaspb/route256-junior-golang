package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

func (c *CliService) list(customerID models.IDType, lastN uint) error {
	orders, err := c.srvc.GetCustomerOrders(c.ctx, customerID, lastN)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if len(orders) == 0 {
		fmt.Println("No orders")
		return nil
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
	return nil
}

//func init() {
//	rootCli.AddCommand(listCli)
//
//	listCli.AddCommand(&cobra.Command{
//		Use:   "help",
//		Short: "Help about command",
//		Long:  `Help about command`,
//		Run: func(cmd *cobra.Command, args []string) {
//			if err := listCli.Help(); err != nil {
//				fmt.Println(err.Error())
//			}
//		},
//	})
//
//	listCli.Flags().UintP("n", "n", 0, "N last orders count")
//}

func (c *CliService) initListCmd(rootCli *cobra.Command) {
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

			if err := c.list(customerID, n); err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	rootCli.AddCommand(listCli)

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
