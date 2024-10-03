package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
)

func GiveHandler(ctx context.Context, buffer *bytes.Buffer, service *srvc.Service, customerID models.IDType, orderIDs []models.IDType) error {
	orders, err := service.GiveOrderToCustomer(ctx, orderIDs, customerID)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if len(orders) == 0 {
		fmt.Fprintln(buffer, "No orders")
		return nil
	}

	maxMsgLen := 0
	maxPackageNameLen := len("Pack")
	for _, order := range orders {
		maxMsgLen = max(maxMsgLen, len(order.Msg))
		maxPackageNameLen = max(maxPackageNameLen, len(order.Package))
	}

	tableTop := fmt.Sprintf(
		"%8s|Give|%-"+strconv.Itoa(maxMsgLen)+"s|%"+strconv.Itoa(maxPackageNameLen)+"s|Cost",
		"ID",
		"Message",
		"Pack",
	)
	fmt.Fprintln(buffer, tableTop)
	for _, order := range orders {
		give := "NO"
		if order.Ok {
			give = "YES"
		}
		tableRow := fmt.Sprintf(
			"%8v|%4s|%-"+strconv.Itoa(maxMsgLen)+"s|%-"+strconv.Itoa(maxPackageNameLen)+"s|%v",
			order.ID, give, order.Msg, order.Package, order.Cost)
		fmt.Fprintln(buffer, tableRow)
	}
	return nil
}

func getGiveCmd(service *srvc.Service) *cobra.Command {
	var giveCli = &cobra.Command{
		Use:     "give",
		Short:   "Give orders by their ids from PVZ to customer",
		Long:    `Give orders by their ids from PVZ to customer`,
		Example: "give <userID> <orderID_1> ... <orderID_N>",
		Run: func(cmd *cobra.Command, args []string) {
			if startTime, err := getStartTimeInCmd(cmd); err != nil {
				if !errors.Is(err, ErrorNoStartTimeInCMD) {
					fmt.Println(err.Error())
					return
				}
			} else {
				service.SetStartTime(startTime)
			}

			if len(args) < 2 {
				fmt.Println("IDs count is less then 2, check 'GiveHandler --help'")
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

			customerID := ids[0]
			orderIDs := ids[1:]

			var buffer bytes.Buffer
			if err := GiveHandler(cmd.Context(), &buffer, service, customerID, orderIDs); err != nil {
				fmt.Println(err.Error())
			}
			fmt.Print(buffer.String())
		},
	}

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

	return giveCli
}
