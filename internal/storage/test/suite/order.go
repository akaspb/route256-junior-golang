package suite

import (
	"gitlab.ozon.dev/siralexpeter/Homework/internal/models"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/storage"
	"gitlab.ozon.dev/siralexpeter/Homework/test/helpers"
	"time"

	"github.com/stretchr/testify/suite"
)

const testJsonPath = "test_storage.json"

type ItemSuite struct {
	suite.Suite
	storage storage.Storage
}

func (s *ItemSuite) SetupSuite() {
	if err := helpers.DeleteFile(testJsonPath); err != nil {
		s.Fail(err.Error())
	}
	if err := helpers.CopyFile("test_storages/storage.json", testJsonPath); err != nil {
		s.Fail(err.Error())
	}

	testStorage, err := storage.InitJsonStorage(testJsonPath)
	s.storage = testStorage
	s.Require().NoError(err)
}

func (s *ItemSuite) TestSetOrder() {
	now := time.Now()
	day := 24 * time.Hour
	month := 30 * day

	order := models.Order{
		ID:         0,
		CustomerID: 0,
		Expiry:     now.Add(month),
		Weight:     1,
		Cost:       1,
		Pack:       nil,
		Status: models.Status{
			Value:     models.StatusToStorage,
			ChangedAt: now,
		},
	}

	err := s.storage.SetOrder(order)
	s.Require().NoError(err)
}

func (s *ItemSuite) TestGetOrder() {
	_, err := s.storage.GetOrder(3)
	s.Require().NoError(err)

	_, err = s.storage.GetOrder(7)
	s.Require().Error(err)
}

func (s *ItemSuite) TestGetOrderNotInStorage() {
	_, err := s.storage.GetOrder(7)
	s.Require().Error(err)
}

//
//func (s *ItemSuite) TestAddItemSuccess() {
//	ctx := context.Background()
//
//	err := s.itemHandler.AddItem(ctx, 10, domain.Item{
//		SKU:   1076963,
//		Count: 10,
//	})
//
//	s.Require().NoError(err)
//}
//
//func (s *ItemSuite) TestAddItemFailed() {
//	ctx := context.Background()
//
//	err := s.itemHandler.AddItem(ctx, 10, domain.Item{
//		SKU:   1,
//		Count: 10,
//	})
//
//	s.Require().Error(err)
//}
