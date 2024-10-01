package cli

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
)

var (
	ErrorNoStartTimeInCMD = errors.New("no start time flag in cmd")
)

type CliService struct {
	ctx          context.Context
	orderStorage storage.Facade
	packService  *packaging.Packaging
	srvc         *srvc.Service
	rootCli      *cobra.Command
}

func NewCliService(ctx context.Context, orderStorage storage.Facade, packService *packaging.Packaging, service *srvc.Service) *CliService {
	c := &CliService{
		ctx:          ctx,
		orderStorage: orderStorage,
		packService:  packService,
		srvc:         service,
	}

	c.rootCli = getRootCli()
	c.rootCli.AddCommand(getGiveCmd(ctx, service))
	c.rootCli.AddCommand(getInterCmd(service, c.rootCli))
	c.rootCli.AddCommand(getListCmd(ctx, service))
	c.rootCli.AddCommand(getNowCmd(service))
	c.rootCli.AddCommand(getReceiveCmd(ctx, service, packService))
	c.rootCli.AddCommand(getRemoveCmd(ctx, service))
	c.rootCli.AddCommand(getReturnCmd(ctx, service))
	c.rootCli.AddCommand(getReturnsCmd(ctx, service))

	return c
}

//func (c *CliService) updateTimeInServiceInCmd(cmd *cobra.Command) error {
//	if cmd.Flags().Changed("start") {
//		startStr, err := cmd.Flags().GetString("start")
//		if err != nil {
//			return err
//		}
//
//		startTime, err := time.Parse("02.01.2006", startStr)
//		if err != nil {
//			return err
//		}
//
//		c.srvc.SetStartTime(startTime)
//	}
//
//	return nil
//}

func getStartTimeInCmd(cmd *cobra.Command) (time.Time, error) {
	if !cmd.Flags().Changed("start") {
		return time.Time{}, ErrorNoStartTimeInCMD
	}

	startStr, err := cmd.Flags().GetString("start")
	if err != nil {
		return time.Time{}, err
	}

	startTime, err := time.Parse("02.01.2006", startStr)
	if err != nil {
		return time.Time{}, err
	}

	return startTime, nil
}

func getToday(service *srvc.Service) string {
	return service.GetCurrentTime().Truncate(24 * time.Hour).Format("02.01.2006")
}

func (c *CliService) Execute() error {
	var err error

	err = c.rootCli.Execute()
	if err != nil {
		return err
	}

	return nil
}
