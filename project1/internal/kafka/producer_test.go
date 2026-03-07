package kafka

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/dhruv/kafka-microservices/internal/models"
)

func TestProducer_PublishEvent(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Use sarama mocks to avoid needing a real Kafka broker
	mockProducer := mocks.NewSyncProducer(t, config)

	// Expect a message to be published successfully
	mockProducer.ExpectSendMessageAndSucceed()

	// We can inject the mock into our Producer wrapper.
	// We have to mock the wrapper initialization or create it manually:
	producerWrapper := &Producer{
		syncProducer: mockProducer,
	}

	order := models.Order{
		ID:         "test-id",
		Item:       "Phone",
		LegacyQty:  2,
		LegacyPrice: 500.0,
	}

	err := producerWrapper.PublishEvent("orders.created", order)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := producerWrapper.Close(); err != nil {
		t.Fatalf("Error closing mock producer: %v", err)
	}
}
