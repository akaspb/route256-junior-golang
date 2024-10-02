package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
)

type ReceiveHandlerDTO struct {
	Ctx         context.Context
	Buffer      *bytes.Buffer
	Service     *srvc.Service
	OrderID     models.IDType
	OrderCost   models.CostType
	OrderWeight models.WeightType
	OrderExpiry time.Time
	PackName    string
	CustomerID  models.IDType
}

func ReceiveHandler(
	receiveHandlerDTO ReceiveHandlerDTO,
) error {
	ctx := receiveHandlerDTO.Ctx
	buffer := receiveHandlerDTO.Buffer
	service := receiveHandlerDTO.Service
	orderID := receiveHandlerDTO.OrderID
	orderCost := receiveHandlerDTO.OrderCost
	orderWeight := receiveHandlerDTO.OrderWeight
	orderExpiry := receiveHandlerDTO.OrderExpiry
	packName := receiveHandlerDTO.PackName
	customerID := receiveHandlerDTO.CustomerID

	var packPtr *models.Pack
	if packName != "" {
		pack, err := service.Packaging.GetPackagingByName(packName)
		if err != nil {
			return err
		}

		packPtr = &pack
	}

	if err := service.AcceptOrderFromCourier(ctx, srvc.AcceptOrderDTO{
		OrderID:     orderID,
		OrderCost:   orderCost,
		OderWeight:  orderWeight,
		CustomerID:  customerID,
		Pack:        packPtr,
		OrderExpiry: orderExpiry,
	}); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	fmt.Fprintln(buffer, "success: take order for storage in PVZ")
	return nil
}

func getReceiveCmd(ctx context.Context, service *srvc.Service, packService *packaging.Packaging) *cobra.Command {
	var receiveCli = &cobra.Command{
		Use:     "receive",
		Short:   "Receive order from courier to PVZ",
		Long:    `Receive order from courier to PVZ`,
		Example: "receive -o=<orderID> -m=<order cost> -w=<order weight> -c=<customerID> -e=<expiry in DD.MM.YYYY format> [-p=<packaging name>]",
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

			var buffer bytes.Buffer
			if err := ReceiveHandler(
				ReceiveHandlerDTO{
					Ctx:         ctx,
					Buffer:      &buffer,
					Service:     service,
					OrderID:     orderID,
					OrderCost:   cost,
					OrderWeight: weight,
					OrderExpiry: orderExpiry,
					PackName:    packName,
					CustomerID:  customerID,
				},
			); err != nil {
				fmt.Println(err.Error())
			}
			fmt.Print(buffer.String())
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

	packs := packService.GetAllPacks()
	packagingNames := make([]string, 0, len(packs))
	for _, pack := range packs {
		packagingNames = append(packagingNames, fmt.Sprintf("\n\t - %s", pack.Name))
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

	return receiveCli
}
