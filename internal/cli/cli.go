package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/status"
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
	c.rootCli.AddCommand(getInterCmd(c.rootCli))
	c.rootCli.AddCommand(getListCmd(client))
	c.rootCli.AddCommand(getReceiveCmd(client))
	c.rootCli.AddCommand(getRemoveCmd(client))
	c.rootCli.AddCommand(getReturnCmd(client))
	c.rootCli.AddCommand(getReturnsCmd(client))
	c.rootCli.AddCommand(getThreadsCmd(client))

	return c
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
	errStatus, ok := status.FromError(err)
	if !ok {
		fmt.Println("handleResponseError function should be used with status errors only")
		return
	}

	fmt.Printf("%v: %v\n", errStatus.Code(), errStatus.Message())
}
