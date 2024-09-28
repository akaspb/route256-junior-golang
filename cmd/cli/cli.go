package cli

import (
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
)

var CliServiceGlobal CliService

type CliService struct {
	orderStorage storage.Facade
	packService  packaging.Packaging
	srvc         srvc.Service
}

func NewCliService(orderStorage storage.Facade, packService packaging.Packaging, service srvc.Service) *CliService {
	//var (
	//	in        = bufio.NewReader(os.Stdin)
	//	todayTime = time.Now().Truncate(24 * time.Hour)
	//	todayStr  = todayTime.Format("02.01.2006")
	//	packs     []models.Pack
	//)

	return &CliService{
		orderStorage: orderStorage,
		packService:  packService,
		srvc:         service,
	}
}

func (cs *CliService) updateTimeInService(currTime time.Time) error {
	cs.srvc.StartTime = currTime

	return nil
}

func (cs *CliService) updateTimeInServiceInCmd(cmd *cobra.Command) error {
	if cmd.Flags().Changed("start") {
		startStr, err := cmd.Flags().GetString("start")
		if err != nil {
			return err
		}

		startTime, err := time.Parse("02.01.2006", startStr)
		if err != nil {
			return err
		}

		return cs.updateTimeInService(startTime)
	}

	return nil
}

func (cs *CliService) Execute() error {
	var err error

	err = rootCli.Execute()
	if err != nil {
		return err
	}

	return nil
}
