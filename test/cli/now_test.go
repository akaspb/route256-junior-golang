package test_cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/cli"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/packaging"
	srvc "gitlab.ozon.dev/siralexpeter/Homework/internal/service"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers/storage"
)

func TestNowCmd(t *testing.T) {
	testStorage := storage.NewStorage()

	packService, err := packaging.NewPackaging()
	if err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	startTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	correctResult := strings.TrimSpace(startTime.Truncate(24 * time.Hour).Format("02.01.2006"))

	service := srvc.NewService(testStorage, packService, startTime, time.Now())

	var buffer bytes.Buffer
	err = cli.NowHandler(&buffer, service)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buffer.String()
	output = strings.TrimSpace(output)

	if output != correctResult {
		t.Errorf("wrong output: %v", output)
	}
}
