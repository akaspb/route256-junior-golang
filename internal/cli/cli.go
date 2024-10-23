package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

func AddInterToRoot(root *cobra.Command) {
	root.AddCommand(getInterCmd(root))
}

type CliService struct {
	rootCli *cobra.Command
}

func NewCliService(root *cobra.Command, cmds ...*cobra.Command) *CliService {
	c := &CliService{
		rootCli: root,
	}

	for _, cmd := range cmds {
		c.rootCli.AddCommand(cmd)
	}

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
