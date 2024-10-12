package cli

import (
	"context"
	"errors"
	"fmt"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"google.golang.org/grpc/status"
	"time"

	"github.com/spf13/cobra"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
)

var (
	ErrorNoStartTimeInCMD = errors.New("no start time flag in cmd")
)

type CliService struct {
	client  pvz_service.PvzServiceClient
	rootCli *cobra.Command
}

func NewCliService(client pvz_service.PvzServiceClient) *CliService {
	c := &CliService{
		client: client,
	}

	c.rootCli = getRootCli()
	c.rootCli.AddCommand(getGiveCmd(client))
	c.rootCli.AddCommand(getInterCmd(service, c.rootCli))
	c.rootCli.AddCommand(getListCmd(service))
	c.rootCli.AddCommand(getReceiveCmd(service, packService))
	c.rootCli.AddCommand(getRemoveCmd(service))
	c.rootCli.AddCommand(getReturnCmd(service))
	c.rootCli.AddCommand(getReturnsCmd(service))
	c.rootCli.AddCommand(getThreadsCmd(service))

	return c
}

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

func (c *CliService) Execute(ctx context.Context) error {
	var err error

	err = c.rootCli.ExecuteContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func handleResponseError(err error) {
	code := status.Code(err)
	fmt.Printf("%v: %v\n", code, err)
}
