package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/internal/pvz-service/v1"
)

func GetRemoveCmd(client pvz_service.PvzServiceClient) *cobra.Command {
	var removeCli = &cobra.Command{
		Use:     "remove",
		Short:   "Return order from PVZ to courier",
		Long:    `Return order from PVZ to courier`,
		Example: "remove <orderID>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("orderID is not defined, check 'remove --help'")
				return
			}

			orderID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			request := &pvz_service.RemoveOrderRequest{
				OrderId: orderID,
			}

			_, err = client.RemoveOrder(cmd.Context(), request)
			if err != nil {
				handleResponseError(err)
				return
			}

			fmt.Println("success")
		},
	}

	removeCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := removeCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})

	return removeCli
}
