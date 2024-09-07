package cli

import (
	"bufio"
	"github.com/spf13/cobra"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"os"
	"time"
)

const (
	jsonPath = "internal/storage/storage.json"
)

var in = bufio.NewReader(os.Stdin)
var todayTime = time.Now().Truncate(24 * time.Hour)
var todayStr = todayTime.Format("02.01.2006")
var cliService = NewCliService()

type CliService struct {
	srvc *srvc.Service
}

func NewCliService() *CliService {
	return &CliService{}
}

func (cs *CliService) getService(currTime time.Time) error {
	orderStorage, err := storage.InitJsonStorage(jsonPath)
	if err != nil {
		return err
	}

	cs.srvc = srvc.NewService(orderStorage, currTime, time.Now().Truncate(24*time.Hour))
	return nil
}

func (cs *CliService) getServiceInCommand(cmd *cobra.Command) error {
	if cmd.Flags().Changed("start") {
		startStr, err := cmd.Flags().GetString("start")
		if err != nil {
			return err
		}

		startTime, err := time.Parse("02.01.2006", startStr)
		if err != nil {
			return err
		}

		return cs.getService(startTime)
	}

	if cs.srvc == nil {
		return cs.getService(todayTime)
	}

	return nil
}
