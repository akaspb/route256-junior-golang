package cli

import (
	"bytes"
	"context"
	"fmt"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

func ListHandler(ctx context.Context, buffer *bytes.Buffer, service *srvc.Service, customerID models.IDType, lastN uint) error {
	orders, err := service.GetCustomerOrders(ctx, customerID, lastN)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if len(orders) == 0 {
		fmt.Fprintln(buffer, "No orders")
		return nil
	}

	maxPackageNameLen := len("Pack")
	for _, order := range orders {
		maxPackageNameLen = max(maxPackageNameLen, len(order.Package))
	}

	tableTop := fmt.Sprintf("%8s|    Expiry|Expired|%"+strconv.Itoa(maxPackageNameLen)+"s|Cost", "ID", "Pack")
	fmt.Fprintln(buffer, tableTop)

	for _, order := range orders {
		expired := "NO"
		if order.Expired {
			expired = "YES"
		}

		tableRow := fmt.Sprintf(
			"%8v|%-10s|%7s|%"+strconv.Itoa(maxPackageNameLen)+"s|%v",
			order.ID, order.Expiry.Format("02.01.2006"), expired, order.Package, order.Cost,
		)
		fmt.Fprintln(buffer, tableRow)
	}
	return nil
}

func getListCmd(ctx context.Context, service *srvc.Service) *cobra.Command {
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

			var buffer bytes.Buffer
			if err := ListHandler(ctx, &buffer, service, customerID, n); err != nil {
				fmt.Println(err.Error())
			}
			fmt.Print(buffer.String())
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
