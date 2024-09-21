package cli

import (
	"bufio"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	"os"
	"time"

	"github.com/spf13/cobra"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
)

const (
	relativeJsonPath = /*"D:/Go/Ozon/Homework/" +*/ "internal/storage/storage.json"
)

var (
	in         = bufio.NewReader(os.Stdin)
	todayTime  = time.Now().Truncate(24 * time.Hour)
	todayStr   = todayTime.Format("02.01.2006")
	cliService = NewCliService(relativeJsonPath)
	packs      []models.Pack
)

func init() {
	//absoluteJsonPath = "D:/Go/Ozon/Homework/" + relativeJsonPath

	packet := packaging.NewPack("packet", 5, 10)
	box := packaging.NewPack("box", 20, 30)
	wrap := packaging.NewPack("wrap", 1, packaging.AnyWeight)

	packs = []models.Pack{
		packet,
		box,
		wrap,
	}
}

type CliService struct {
	jsonPath string
	srvc     *srvc.Service
}

func NewCliService(jsonPath string) *CliService {
	return &CliService{jsonPath: jsonPath}
}

func (cs *CliService) getService(currTime time.Time) error {
	orderStorage, err := storage.InitJsonStorage(cs.jsonPath)
	if err != nil {
		return err
	}

	packService, err := packaging.NewPackaging(packs...)
	if err != nil {
		return err
	}

	cs.srvc = srvc.NewService(orderStorage, packService, currTime, time.Now().Truncate(24*time.Hour))
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

func Execute() error {
	var err error

	err = rootCli.Execute()
	if err != nil {
		return err
	}

	return nil
}
