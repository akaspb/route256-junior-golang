package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcServerHost = "localhost:7001"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer cancel()
	serviceStopped := make(chan struct{})

	conn, err := grpc.NewClient(grpcServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer conn.Close()

	pvzClient := pvz_service.NewPvzServiceClient(conn)
	pvzCmd := cli.NewCliService(pvzClient)

	go func() {
		if err := pvzCmd.Execute(ctx); err != nil {
			fmt.Printf("error during pvz command execution: %v\n", err)
		}
		close(serviceStopped)
	}()

	select {
	case <-ctx.Done():
		<-serviceStopped
	case <-serviceStopped:
	}
}
