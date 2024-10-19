package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	pvz_service "gitlab.ozon.dev/siralexpeter/Homework/internal/pvz-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer cancel()
	serviceStopped := make(chan struct{})

	if err := initConfig(); err != nil {
		fmt.Printf("error initializing configs: %s\n", err.Error())
		return
	}

	conn, err := grpc.NewClient(
		viper.GetString("grpc.host"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer conn.Close()

	pvzClient := pvz_service.NewPvzServiceClient(conn)
	packService, err := packaging.NewPackaging()
	if err != nil {
		fmt.Printf("error during starting packing service: %v\n", err)
		return
	}

	pvzCmd := cli.NewCliService(pvzClient, packService)

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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
