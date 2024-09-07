package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"time"
)

var receiveCli = &cobra.Command{
	Use:     "receive",
	Short:   "Receive order from courier to PVZ",
	Long:    `Receive order from courier to PVZ`,
	Example: "receive -o=<orderID> -c=<customerID> -e=\"<expiry in DD.MM.YYYY format>\"",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cliService.getServiceInCommand(cmd); err != nil {
			fmt.Println(err.Error())
			return
		}

		var orderID, customerID models.IDType
		var orderExpiry time.Time

		if cmd.Flags().Changed("order") {
			orderIDint64, err := cmd.Flags().GetInt64("order")
			orderID = models.IDType(orderIDint64)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("order flag is not defined")
			return
		}

		if cmd.Flags().Changed("customer") {
			customerIDint64, err := cmd.Flags().GetInt64("customer")
			customerID = models.IDType(customerIDint64)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("customer flag is not defined")
			return
		}

		if cmd.Flags().Changed("expiry") {
			expiryStr, err := cmd.Flags().GetString("expiry")
			if err != nil {
				fmt.Println(err.Error())
			}

			orderExpiry, err = time.Parse("02.01.2006", expiryStr)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("expiry flag is not defined")
			return
		}

		if err := cliService.srvc.AcceptOrderFromCourier(orderID, customerID, orderExpiry); err != nil {
			fmt.Println("error:", err.Error())
		} else {
			fmt.Println("success: take order for storage in PVZ")
		}
	},
}

func init() {
	rootCli.AddCommand(receiveCli)

	receiveCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := receiveCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})

	receiveCli.Flags().Int64P("order", "o", 0, "unique order ID")
	receiveCli.Flags().Int64P("customer", "c", 0, "unique customer ID")
	receiveCli.Flags().StringP("expiry", "e", "", "expiry time in format DD.MM.YYYY")
}
