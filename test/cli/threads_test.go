package test_cli

import (
	"errors"
	"runtime"
	"testing"
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers/storage"
)

func TestThreadsCmd(t *testing.T) {
	testStorage := storage.NewStorage()

	nowTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)

	packService, err := packaging.NewPackaging()
	if err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	service, err := srvc.NewService(testStorage, packService, nowTime, time.Now(), workerCount)
	if err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	tests := []struct {
		newWorkerCount int
		err            error
	}{
		{
			newWorkerCount: 0,
			err:            srvc.ErrorInvalidWorkerCount,
		},
		{
			newWorkerCount: runtime.GOMAXPROCS(0) + 1,
			err:            srvc.ErrorInvalidWorkerCount,
		},
		{
			newWorkerCount: workerCount,
			err:            nil,
		},
	}

	for _, tc := range tests {
		err = cli.ThreadsHandler(service, tc.newWorkerCount)

		if err != nil {
			if !errors.Is(err, tc.err) {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}
