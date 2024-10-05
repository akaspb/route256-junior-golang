package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage/postgres"
)

const (
	psqlDSN     = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"
	workerCount = 2
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx) // syscall.SIGINT, syscall.SIGTERM
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

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
	service, err := srvc.NewService(orderStorage, packService, now, now, workerCount)
	if err != nil {
		fmt.Printf("error in main func: %v\n", err)
		return
	}

	cliService := cli.NewCliService(orderStorage, packService, service)

	err = cliService.Execute(ctx)
	if err != nil {
		fmt.Printf("error in main func: %v\n", err)
		return
	}

	<-sigCh
	fmt.Println("$$$ program finish")
}

func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgStorage(txManager)

	return storage.NewStorageFacade(txManager, pgRepository)
}
