package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func GetRootCli() *cobra.Command {
	return &cobra.Command{
		Use:   "pvz",
		Short: "pvz is program for implementing the interaction of the PVZ manager with the courier and the customer",
		Long:  `pvz is program for implementing the interaction of the PVZ manager with the courier and the customer`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("to start interactive mode use 'pvz inter'")
		},
	}
}
