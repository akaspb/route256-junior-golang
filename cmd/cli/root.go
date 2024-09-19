package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCli = &cobra.Command{
	Use:   "pvz",
	Short: "pvz is program for implementing the interaction of the PVZ manager with the courier and the customer",
	Long:  `pvz is program for implementing the interaction of the PVZ manager with the courier and the customer`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("to start interactive mode use 'pvz inter'")
	},
}

func Execute() {

	if err := rootCli.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
