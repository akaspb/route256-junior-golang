package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

func RemoveHandler(ctx context.Context, buffer *bytes.Buffer, service *srvc.Service, orderID models.IDType) error {
	if err := service.ReturnOrder(ctx, orderID); err != nil {
		return fmt.Errorf("error: %v\n", err)
	}

	fmt.Fprintln(buffer, "success: order can be given to courier for return")

	return nil
}

func getRemoveCmd(ctx context.Context, service *srvc.Service) *cobra.Command {
	var removeCli = &cobra.Command{
		Use:     "remove",
		Short:   "Return order from PVZ to courier",
		Long:    `Return order from PVZ to courier`,
		Example: "remove <orderID>",
		Run: func(cmd *cobra.Command, args []string) {
			if startTime, err := getStartTimeInCmd(cmd); err != nil {
				if !errors.Is(err, ErrorNoStartTimeInCMD) {
					fmt.Println(err.Error())
					return
				}
			} else {
				service.SetStartTime(startTime)
			}

			if len(args) < 1 {
				fmt.Println("orderID is not defined, check 'remove --help'")
				return
			}

			orderIDint64, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			orderID := models.IDType(orderIDint64)

			var buffer bytes.Buffer
			err = RemoveHandler(ctx, &buffer, service, orderID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Print(buffer.String())
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
