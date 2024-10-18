package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	event_factory "gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger/factory"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger/kafka"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/middleware"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/server"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage/postgres"
	desc "gitlab.ozon.dev/siralexpeter/Homework/pkg/pvz-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	if err := initConfig(); err != nil {
		fmt.Printf("error initializing configs: %s\n", err.Error())
		return
	}

	pool, err := pgxpool.Connect(ctx, getPostgresDSN())
	if err != nil {
		fmt.Printf("error in main func: %v\n", err)
		return
	}
	defer pool.Close()

	orderStorage := newStorageFacade(pool)

	packService, err := packaging.NewPackaging()
	if err != nil {
		fmt.Printf("error in main func: %v\n", err)
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	pvzService, err := service.NewService(orderStorage, packService, now, now, viper.GetInt("program.worker_count"))
	if err != nil {
		fmt.Printf("error in main func: %v\n", err)
		return
	}
	defer pvzService.Close()

	kafkaProducer, err := initProducer(kafka.Config{
		Brokers: []string{viper.GetString("kafka.broker")},
	})
	if err != nil {
		fmt.Printf("failed to init kafka producer: %v", err)
		return
	}

	kafkaEventLogger := kafka.NewTopicSender(kafkaProducer, viper.GetString("kafka.topic"))
	eventFactory := event_factory.NewDefaultFactory(1)
	pvzServer := server.NewImplementation(pvzService, kafkaEventLogger, eventFactory)

	lis, err := net.Listen("tcp", viper.GetString("grpc.host"))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.Logging),
	)
	reflection.Register(grpcServer)
	desc.RegisterPvzServiceServer(grpcServer, pvzServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v", err)
			return
		}
	}()

	select {
	case <-ctx.Done():
	}

	kafkaProducer.AsyncClose()
	<-kafkaProducer.Successes()
	<-kafkaProducer.Errors()

	grpcServer.GracefulStop()
}

func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgStorage(txManager)

	return storage.NewStorageFacade(txManager, pgRepository)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func getPostgresDSN() string {
	//postgres://username:password@host:port/database
	//"postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.postgres"),
		viper.GetString("db.dbname"),
		viper.GetString("db.sslmode"),
	)
}
