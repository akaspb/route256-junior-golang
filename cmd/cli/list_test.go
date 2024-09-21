package cli

import (
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers"
	"testing"
)

func TestList(t *testing.T) {
	jsonFileCopy := "/TestList_storage.json"
	if err := helpers.DeleteFile(jsonFileCopy); err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}
	if err := helpers.CopyFile("./test_storages/list_test_storage.json", jsonFileCopy); err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	cliService := NewCliService(jsonFileCopy)
	if err := cliService.getService(todayTime); err != nil {
		t.Fatalf("unexpected error before test: %v", err)
	}

	//if err := cliService.list(0, 100); err != nil {
	//	t.Errorf("unexpected error: %v", err)
	//}
}
