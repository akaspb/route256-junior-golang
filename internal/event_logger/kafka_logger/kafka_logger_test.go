package kafka_logger

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"testing"
	"time"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger"
	"gitlab.ozon.dev/siralexpeter/Homework/internal/kafka/mocks"
)

const testTopic = "test.events-log"

func TestKafkaLogger_Send(t *testing.T) {
	currentTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

	asyncProducerMock := mocks.NewAsyncProducerMock(t)
	kafkaLogger := NewKafkaLogger(asyncProducerMock, testTopic)

	tests := []event_logger.Event{
		{
			ID:        1,
			Type:      "Test1",
			Details:   "Test1",
			Timestamp: currentTime,
		},
		{
			ID:   2,
			Type: "Test2",
			Details: struct {
				A string `json:"a"`
				B string `json:"b"`
			}{
				A: "a letter",
				B: "b letter",
			},
			Timestamp: currentTime,
		},
	}

	for _, event := range tests {
		t.Logf("[info] event for test: %v", event)

		correctBytes, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("unexpected error during event marshalling: %v", err)
		}

		producerInputChan := make(chan *sarama.ProducerMessage, 1)
		asyncProducerMock.InputMock.Expect().Return(producerInputChan)

		err = kafkaLogger.Send(event)
		if err != nil {
			switch err.(type) {
			case *json.UnsupportedTypeError:
				t.Errorf("UnsupportedTypeError error during sending event: %v", err)
			case *json.UnsupportedValueError:
				t.Errorf("UnsupportedValueError error during sending event: %v", err)
			default:
				t.Errorf("error during sending event: %v", err)
			}
		}

		producerMessage := <-producerInputChan
		bytesAfterSending, err := producerMessage.Value.Encode()
		if err != nil {
			t.Fatalf("unexpected error during getting bytes from message value: %v", err)
		}

		if !compareBytes(correctBytes, bytesAfterSending) {
			t.Errorf("event bytes before and after sending are not the same")
		}

		close(producerInputChan)
	}
}

func compareBytes(bytes1, bytes2 []byte) bool {
	if len(bytes1) != len(bytes2) {
		return false
	}

	for i := 0; i < len(bytes1); i++ {
		if bytes1[i] != bytes2[i] {
			return false
		}
	}

	return true
}
