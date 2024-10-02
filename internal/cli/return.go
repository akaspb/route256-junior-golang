package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
)

func ReturnHandler(ctx context.Context, buffer *bytes.Buffer, service *srvc.Service, customerID, orderID models.IDType) error {
	if err := service.ReturnOrderFromCustomer(ctx, customerID, orderID); err != nil {
		return err
	}

	fmt.Fprintln(buffer, "success: take order from customer to store it in PVZ")

	return nil
}

func getReturnCmd(ctx context.Context, service *srvc.Service) *cobra.Command {
	var returnCli = &cobra.Command{
		Use:     "return",
		Short:   "Get order from customer to return",
		Long:    `Get order from customer to return`,
		Example: "return -o=<orderID> -c=<customerID>",
		Run: func(cmd *cobra.Command, args []string) {
			if startTime, err := getStartTimeInCmd(cmd); err != nil {
				if !errors.Is(err, ErrorNoStartTimeInCMD) {
					fmt.Println(err.Error())
					return
				}
			} else {
				service.SetStartTime(startTime)
			}

			var orderID, customerID models.IDType

			if !cmd.Flags().Changed("order") {
				fmt.Println("order flag is not defined, check 'return --help'")
				return
			}
			orderIDint64, err := cmd.Flags().GetInt64("order")
			orderID = models.IDType(orderIDint64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !cmd.Flags().Changed("customer") {
				fmt.Println("customer flag is not defined, check 'return --help'")
				return
			}
			customerIDint64, err := cmd.Flags().GetInt64("customer")
			customerID = models.IDType(customerIDint64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var buffer bytes.Buffer
			err = ReturnHandler(ctx, &buffer, service, customerID, orderID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Print(buffer.String())
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
