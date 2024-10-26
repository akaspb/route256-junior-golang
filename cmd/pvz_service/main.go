package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	event_factory "gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger/factory"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger/kafka_logger"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/http"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/kafka"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/metrics"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/middleware"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	desc "gitlab.ozon.dev/siralexpeter/Homework/internal/pvz-service/v1"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/server"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage/in_memory_cache"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage/postgres"
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
		fmt.Printf("error during creating a new db pool: %v\n", err)
		return
	}
	defer pool.Close()

	orderStorage := newStorageFacade(pool)

	packService, err := packaging.NewPackaging()
	if err != nil {
		fmt.Printf("error during starting packing service: %v\n", err)
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	pvzService := service.NewService(
		in_memory_cache.NewInMemoryCache(
			orderStorage,
			viper.GetDuration("service.cache_ttl"),
		),
		packService,
		now,
		now,
	)

	metrics.Init()

	httpServer := http.NewServer(viper.GetString("service.http_address"))
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Printf("http server ListenAndServe: %v\n", err)
		}
	}()

	kafkaProducer, err := initProducer(kafka.Config{
		Brokers: viper.GetStringSlice("kafka_logger.brokers"),
	})
	if err != nil {
		fmt.Printf("failed to init kafka producer: %v", err)
		return
	}

	pvzServer := server.NewImplementation(pvzService)

	lis, err := net.Listen("tcp", viper.GetString("grpc.host"))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}

	kafkaEventLogger := kafka_logger.NewKafkaLogger(
		kafka.NewProducerWrapper(kafkaProducer),
		viper.GetString("kafka_logger.topic"),
	)
	eventFactory := event_factory.NewDefaultFactory(1)
	remoteLogging := middleware.GetRemoteLogging(
		kafkaEventLogger,
		eventFactory,
		viper.GetStringSlice("kafka_logger.methods"),
	)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.LocalLogging, remoteLogging),
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

	for kafkaErr := range kafkaProducer.Errors() {
		fmt.Printf("kafka error: %v\n", kafkaErr.Err)
	}

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("http handler shutdown: %v\n", err)
	}

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

func initProducer(config kafka.Config) (sarama.AsyncProducer, error) {
	return kafka.NewAsyncProducer(config,
		kafka.WithMaxRetries(5),
		kafka.WithRetryBackoff(10*time.Millisecond),
		kafka.WithProducerFlushMessages(2),
		kafka.WithProducerFlushFrequency(500*time.Millisecond),
	)
}
