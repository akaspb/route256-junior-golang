package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

func (c *CliService) initReturnCmd(rootCli *cobra.Command) {
	var returnCli = &cobra.Command{
		Use:     "return",
		Short:   "Get order from customer to return",
		Long:    `Get order from customer to return`,
		Example: "return -o=<orderID> -c=<customerID>",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.updateTimeInServiceInCmd(cmd); err != nil {
				fmt.Println(err.Error())
				return
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

			if err := c.srvc.ReturnOrderFromCustomer(c.ctx, customerID, orderID); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("success: take order from customer to store it in PVZ")
			}
		},
	}

	rootCli.AddCommand(returnCli)

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
}
