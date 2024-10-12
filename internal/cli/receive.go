package cli

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"
	"time"

	"github.com/spf13/cobra"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
)

func getReceiveCmd(client pvz_service.PvzServiceClient) *cobra.Command {
	var receiveCli = &cobra.Command{
		Use:     "receive",
		Short:   "Receive order from courier to PVZ",
		Long:    `Receive order from courier to PVZ`,
		Example: "receive -o=<orderID> -m=<order cost> -w=<order weight> -c=<customerID> -e=<expiry in DD.MM.YYYY format> [-p=<packaging name>]",
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flags().Changed("order") {
				fmt.Println("order flag is not defined, check 'receive --help'")
				return
			}
			orderID, err := cmd.Flags().GetInt64("order")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !cmd.Flags().Changed("weight") {
				fmt.Println("weight flag is not defined, check 'receive --help'")
				return
			}
			weight, err := cmd.Flags().GetFloat32("weight")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !cmd.Flags().Changed("money") {
				fmt.Println("cost flag is not defined, check 'receive --help'")
				return
			}
			cost, err := cmd.Flags().GetFloat32("money")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !cmd.Flags().Changed("customer") {
				fmt.Println("customer flag is not defined, check 'receive --help'")
				return
			}
			customerID, err := cmd.Flags().GetInt64("customer")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !cmd.Flags().Changed("expiry") {
				fmt.Println("expiry flag is not defined, check 'receive --help'")
				return
			}
			expiryStr, err := cmd.Flags().GetString("expiry")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			expiry, err := time.Parse("02.01.2006", expiryStr)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			packName, err := cmd.Flags().GetString("pack")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var packing *wrapperspb.StringValue
			if packName != "" {
				packing = wrapperspb.String(packName)
			}

			request := &pvz_service.ReceiveOrderRequest{
				Id:         orderID,
				CustomerId: customerID,
				Expiry:     timestamppb.New(expiry),
				Weight:     weight,
				Cost:       cost,
				Packing:    packing,
			}

			_, err = client.ReceiveOrder(cmd.Context(), request)
			if err != nil {
				handleResponseError(err)
				return
			}

			fmt.Println("success")
		},
	}

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

	packNames := []string{"packet", "box", "wrap"}
	packs := make([]string, len(packNames))
	for _, packName := range packNames {
		packs = append(packs, fmt.Sprintf("\n\t - %s", packName))
	}

	packagingUsage := strings.Join(packs, "")

	receiveCli.Flags().Int64P("order", "o", 0, "unique order ID")
	receiveCli.Flags().Float32P("money", "m", 0, "order cost")
	receiveCli.Flags().Float32P("weight", "w", 0, "order weight")
	receiveCli.Flags().Int64P("customer", "c", 0, "unique customer ID")
	receiveCli.Flags().StringP("expiry", "e", "", "expiry time in format DD.MM.YYYY")
	receiveCli.Flags().StringP("pack", "p", "", fmt.Sprintf(
		"packaging name: %s",
		packagingUsage,
	))

	return receiveCli
}
