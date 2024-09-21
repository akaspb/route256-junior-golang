package cli

import (
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
)

var receiveCli = &cobra.Command{
	Use:     "receive",
	Short:   "Receive order from courier to PVZ",
	Long:    `Receive order from courier to PVZ`,
	Example: "receive -o=<orderID> -m=<order cost> -w=<order weight> -c=<customerID> -e=<expiry in DD.MM.YYYY format> [-p=<packaging name>]",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cliService.getServiceInCommand(cmd); err != nil {
			fmt.Println(err.Error())
			return
		}

		var orderID, customerID models.IDType
		var orderExpiry time.Time

		if !cmd.Flags().Changed("order") {
			fmt.Println("order flag is not defined, check 'receive --help'")
			return
		}
		orderIDint64, err := cmd.Flags().GetInt64("order")
		orderID = models.IDType(orderIDint64)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if !cmd.Flags().Changed("weight") {
			fmt.Println("weight flag is not defined, check 'receive --help'")
			return
		}
		weightFloat32, err := cmd.Flags().GetFloat32("weight")
		weight := models.WeightType(weightFloat32)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if !cmd.Flags().Changed("money") {
			fmt.Println("cost flag is not defined, check 'receive --help'")
			return
		}
		costFloat32, err := cmd.Flags().GetFloat32("money")
		cost := models.CostType(costFloat32)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if !cmd.Flags().Changed("customer") {
			fmt.Println("customer flag is not defined, check 'receive --help'")
			return
		}
		customerIDint64, err := cmd.Flags().GetInt64("customer")
		customerID = models.IDType(customerIDint64)
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

		orderExpiry, err = time.Parse("02.01.2006", expiryStr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		packName, err := cmd.Flags().GetString("pack")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var packPtr *models.Pack
		if packName != "" {
			pack, err := cliService.srvc.Packaging.GetPackagingByName(packName)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			packPtr = &pack
		}

		if err := cliService.srvc.AcceptOrderFromCourier(service.AcceptOrderDTO{
			OrderID:     orderID,
			OrderCost:   cost,
			OderWeight:  weight,
			CustomerID:  customerID,
			Pack:        packPtr,
			OrderExpiry: orderExpiry,
		}); err != nil {
			fmt.Println("error:", err.Error())
		} else {
			fmt.Println("success: take order for storage in PVZ")
		}
	},
}

func init() {
	RootCli.AddCommand(receiveCli)

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

	packagingNames := make([]string, 0, len(packs))
	for _, packagingName := range packs {
		packagingNames = append(packagingNames, fmt.Sprintf("\n\t - %s", packagingName.Name))
	}

	packagingUsage := strings.Join(packagingNames, "")

	receiveCli.Flags().Int64P("order", "o", 0, "unique order ID")
	receiveCli.Flags().Float32P("money", "m", 0, "order cost")
	receiveCli.Flags().Float32P("weight", "w", 0, "order weight")
	receiveCli.Flags().Int64P("customer", "c", 0, "unique customer ID")
	receiveCli.Flags().StringP("expiry", "e", "", "expiry time in format DD.MM.YYYY")
	receiveCli.Flags().StringP("pack", "p", "", fmt.Sprintf(
		"packaging name: %s",
		packagingUsage,
	))
}
