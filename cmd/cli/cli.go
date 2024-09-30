package cli

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
)

type CliService struct {
	ctx          context.Context
	orderStorage storage.Facade
	packService  packaging.Packaging
	srvc         srvc.Service
	rootCli      *cobra.Command
}

func NewCliService(ctx context.Context, orderStorage storage.Facade, packService packaging.Packaging, service srvc.Service) *CliService {
	cliService := &CliService{
		ctx:          ctx,
		orderStorage: orderStorage,
		packService:  packService,
		srvc:         service,
	}

	cliService.rootCli = getRootCli()
	cliService.initGiveCmd(cliService.rootCli)
	cliService.initInterCmd(cliService.rootCli)
	cliService.initListCmd(cliService.rootCli)
	cliService.initNowCmd(cliService.rootCli)
	cliService.iniReceiveCmd(cliService.rootCli)
	cliService.initRemoveCmd(cliService.rootCli)
	cliService.initReturnCmd(cliService.rootCli)
	cliService.initReturnsCmd(cliService.rootCli)

	return cliService
}

func (c *CliService) updateTimeInServiceInCmd(cmd *cobra.Command) error {
	if cmd.Flags().Changed("start") {
		startStr, err := cmd.Flags().GetString("start")
		if err != nil {
			return err
		}

		startTime, err := time.Parse("02.01.2006", startStr)
		if err != nil {
			return err
		}

		c.srvc.SetStartTime(startTime)
	}

	return nil
}

func (c *CliService) getToday() string {
	return time.Now().Truncate(24 * time.Hour).Format("02.01.2006")
}

func (c *CliService) Execute() error {
	var err error

	err = c.rootCli.Execute()
	if err != nil {
		return err
	}

	return nil
}
