package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	ErrorContextCanceled = errors.New("context was canceled")
)

func getInterCmd(rootCli *cobra.Command) *cobra.Command {
	var interCli = &cobra.Command{
		Use:   "inter",
		Short: "Start program in interactive mode",
		Long:  `Start program in interactive mode`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			for {
				fmt.Print("> ")
				res, err := waitUserInput(ctx)
				if err != nil {
					if !errors.Is(err, io.EOF) {
						fmt.Printf("1\nerror: %v", err)
					}
					return
				}

				res = strings.TrimSpace(res)
				if res == "exit" {
					break
				}

				res = strings.Replace(res, "--help", "help", -1)
				res = strings.Replace(res, "-h", "help", -1)

				argsFromLine := strings.Split(res, " ")
				rootCli.SetArgs(argsFromLine)
				if err := rootCli.Execute(); err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println()
			}
		},
	}

	return interCli
}

func waitUserInput(ctx context.Context) (string, error) {
	group, ctx := errgroup.WithContext(ctx)
	inputChan := make(chan string, 1)
	defer close(inputChan)
	errChan := make(chan error, 1)
	defer close(errChan)

	group.Go(func() error {
		select {
		case <-ctx.Done():
			return ErrorContextCanceled
		case err := <-errChan:
			return err
		}
	})

	go func() {
		in := bufio.NewReader(os.Stdin)
		input, err := in.ReadString('\n')
		if err != nil {
			errChan <- err
		}

		inputChan <- input
		errChan <- nil
	}()

	if err := group.Wait(); err != nil {
		return "", err
	}

	return <-inputChan, nil
}
