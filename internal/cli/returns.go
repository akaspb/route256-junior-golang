package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
)

func ReturnsHandler(
	ctx context.Context,
	buffer *bytes.Buffer,
	service *srvc.Service,
	offset, limit int,
) error {
	returnsChan, err := service.GetReturnsList(ctx, offset, limit)
	if err != nil {
		return fmt.Errorf("error: %v\n", err)

	}

	if len(returnsChan) == 0 {
		fmt.Println()
		return errors.New("No orders")
	}

	fmt.Fprintf(buffer, "%8s|%11s\n", "Order ID", "Customer ID")
	for raw := range returnsChan {
		tableRow := fmt.Sprintf("%8v|%11v", raw.OrderID, raw.CustomerID)
		fmt.Fprintln(buffer, tableRow)
	}

	return nil
}

func getReturnsCmd(service *srvc.Service) *cobra.Command {
	var returnsCli = &cobra.Command{
		Use:     "returns",
		Short:   "Get all orders, which must be given to courier/couriers for return from PVZ",
		Long:    `Get all orders, which must be given to courier/couriers for return from PVZ`,
		Example: "returns -o=<offset> -l=<limit>",
		Run: func(cmd *cobra.Command, args []string) {
			if startTime, err := getStartTimeInCmd(cmd); err != nil {
				if !errors.Is(err, ErrorNoStartTimeInCMD) {
					fmt.Println(err.Error())
					return
				}
			} else {
				service.SetStartTime(startTime)
			}

			var offset, limit int
			var err error

			if !cmd.Flags().Changed("offset") {
				fmt.Println("offset flag is not defined, check 'returns --help'")
				return
			}
			offset, err = cmd.Flags().GetInt("offset")
			if err != nil {
				fmt.Println(err.Error())
			}

			if !cmd.Flags().Changed("limit") {
				fmt.Println("limit flag is not defined, check 'returns --help'")
				return
			}
			limit, err = cmd.Flags().GetInt("limit")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var buffer bytes.Buffer
			err = ReturnsHandler(cmd.Context(), &buffer, service, offset, limit)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Print(buffer.String())
		},
	}

	returnsCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := returnsCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})

	returnsCli.Flags().IntP("limit", "l", 0, "limit of results")
	returnsCli.Flags().IntP("offset", "o", 0, "offset of results, starts from 0")

	return returnsCli
}
