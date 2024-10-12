package cli

import (
	"fmt"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"strconv"

	"github.com/spf13/cobra"
)

//func ThreadsHandler(service *srvc.Service, threadsCount int) error {
//	if err := service.ChangeWorkerCount(threadsCount); err != nil {
//		return fmt.Errorf("error: %w", err)
//	}
//
//	return nil
//}

func getThreadsCmd(client pvz_service.PvzServiceClient) *cobra.Command {
	var threadsCli = &cobra.Command{
		Use:     "threads",
		Short:   "Set max threads count for program tasks processing",
		Long:    `Set max threads count for program tasks processing`,
		Example: "threads <threads count>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("threads count is not given, check 'threads --help'")
				return
			}

			threadsCount, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//if err := ThreadsHandler(service, int(threadsCount)); err != nil {
			//	fmt.Println(err.Error())
			//} else {
			//	fmt.Println("success")
			//}

			request := &pvz_service.ChangeThreadCountRequest{
				ThreadCount: uint32(threadsCount),
			}

			_, err = client.ChangeThreadCount(cmd.Context(), request)
			if err != nil {
				handleResponseError(err)
				return
			}

			fmt.Println("success")
		},
	}

	threadsCli.AddCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about command",
		Long:  `Help about command`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := threadsCli.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
	})

	return threadsCli
}
