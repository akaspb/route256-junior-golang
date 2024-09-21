package main

import (
	"fmt"
	"gitlab.ozon.dev/siralexpeter/Homework/cmd/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Printf("error in main func: %v\n", err)
	}
}
