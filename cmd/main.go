package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage/postgres"
)

const (
	psqlDSN = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, psqlDSN)
	if err != nil {
		fmt.Printf("error in main func: %v\n", err)
		return
	}
	defer pool.Close()

	orderStorage := newStorageFacade(pool)

	packet := packaging.NewPack("packet", 5, 10)
	box := packaging.NewPack("box", 20, 30)
	wrap := packaging.NewPack("wrap", 1, packaging.AnyWeight)

	packService, err := packaging.NewPackaging(packet, box, wrap)
	if err != nil {
		fmt.Printf("error in main func: %v\n", err)
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	service := srvc.NewService(orderStorage, packService, now, now)

	cliService := cli.NewCliService(ctx, orderStorage, packService, service)

	err = cliService.Execute()
	if err != nil {
		fmt.Printf("error in main func: %v\n", err)
		return
	}
}

func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgStorage(txManager)

	return storage.NewStorageFacade(txManager, pgRepository)
}
