package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pvz_cli "gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
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
	pvzCmd := pvz_cli.NewCliService(pvzClient)

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

	//md := metadataParse()
	////ctx = metadata.NewOutgoingContext(ctx, metadata.New(md))
	//ctx = metadata.AppendToOutgoingContext(ctx, md...)

	//var (
	//	resp    proto.Message
	//	respErr error
	//)
	//switch *methodFlag {
	//case "AddContact":
	//	req := &telephone_service.AddContactRequest{}
	//	if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
	//		log.Fatalf("failed to unmarshal data: %v", err)
	//	}
	//	resp, respErr = telephoneServiceClient.AddContact(ctx, req)
	//
	//case "GetContacts":
	//	req := &telephone_service.GetContactsRequest{}
	//	if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
	//		log.Fatalf("failed to unmarshal data: %v", err)
	//	}
	//	resp, respErr = telephoneServiceClient.GetContacts(ctx, req)
	//
	//case "GetContactsHistory":
	//	req := &telephone_service.GetContactsHistoryRequest{}
	//	if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
	//		log.Fatalf("failed to unmarshal data: %v", err)
	//	}
	//	resp, respErr = telephoneServiceClient.GetContactsHistory(ctx, req)
	//
	//case "FindContact":
	//	req := &telephone_service.FindContactRequest{}
	//	if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
	//		log.Fatalf("failed to unmarshal data: %v", err)
	//	}
	//	resp, respErr = telephoneServiceClient.FindContact(ctx, req)
	//
	//case "DeleteContact":
	//	req := &telephone_service.DeleteContactRequest{}
	//	if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
	//		log.Fatalf("failed to unmarshal data: %v", err)
	//	}
	//	resp, respErr = telephoneServiceClient.DeleteContact(ctx, req)
	//
	//default:
	//	log.Fatalf("unknown method: %s", *methodFlag)
	//}

	//data, _ := protojson.Marshal(resp)
	//log.Printf("resp: %v; err: %v\n", string(data), respErr)
	//
	//if status.Code(respErr) == codes.NotFound {
	//	log.Fatalf("no such contact: %s", *methodFlag)
	//}
}

//func metadataParse() []string {
//	md := make(map[string]string)
//	if err := json.Unmarshal([]byte(*metadataFlag), &md); err != nil {
//		log.Fatalf("failed to unmarshal metadata: %v", err)
//	}
//
//	kv := []string{
//		"x-telephone-cli", "21341243123",
//	}
//	for k, v := range md {
//		kv = append(kv, k, v)
//	}
//
//	return kv
//}

//const (
//	psqlDSN     = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"
//	workerCount = 2
//)
//
//func main() {
//	ctx := context.Background()
//	ctx, cancel := context.WithCancel(ctx)
//	defer cancel()
//
//	sigCh := make(chan os.Signal, 1)
//	defer close(sigCh)
//	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
//
//	serviceStopped := make(chan struct{})
//	defer close(serviceStopped)
//
//	pool, err := pgxpool.Connect(ctx, psqlDSN)
//	if err != nil {
//		fmt.Printf("error in main func: %v\n", err)
//		return
//	}
//	defer pool.Close()
//
//	orderStorage := newStorageFacade(pool)
//
//	packService, err := packaging.NewPackaging()
//	if err != nil {
//		fmt.Printf("error in main func: %v\n", err)
//		return
//	}
//
//	now := time.Now().Truncate(24 * time.Hour)
//	service, err := srvc.NewService(orderStorage, packService, now, now, workerCount)
//	if err != nil {
//		fmt.Printf("error in main func: %v\n", err)
//		return
//	}
//
//	cliService := cli.NewCliService(orderStorage, packService, service)
//
//	go func() {
//		if err := cliService.Execute(ctx); err != nil {
//			fmt.Printf("error in main func: %v\n", err)
//		}
//		serviceStopped <- struct{}{}
//	}()
//
//	select {
//	case <-sigCh:
//		cancel()
//		<-serviceStopped
//	case <-serviceStopped:
//		cancel()
//	}
//}
//
//func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
//	txManager := postgres.NewTxManager(pool)
//
//	pgRepository := postgres.NewPgStorage(txManager)
//
//	return storage.NewStorageFacade(txManager, pgRepository)
//}
