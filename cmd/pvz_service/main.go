package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
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

const (
	psqlDSN     = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"
	grpcHost    = "localhost:7001"
	workerCount = 2
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, psqlDSN)
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
	pvzService, err := service.NewService(orderStorage, packService, now, now, workerCount)
	if err != nil {
		fmt.Printf("error in main func: %v\n", err)
		return
	}

	pvzServer := server.NewImplementation(pvzService)

	lis, err := net.Listen("tcp", grpcHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.Logging),
	)
	reflection.Register(grpcServer)
	desc.RegisterPvzServiceServer(grpcServer, pvzServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
	}

	grpcServer.GracefulStop()
}

func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgStorage(txManager)

	return storage.NewStorageFacade(txManager, pgRepository)
}
