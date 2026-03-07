package kafka

import (
	"encoding/json"
	"fmt"
)

type MockProducer struct{}

func NewMockProducer() (*MockProducer, error) {
	return &MockProducer{}, nil
}

func (p *MockProducer) PublishEvent(topic string, event interface{}) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	fmt.Printf("[MOCK] Message published to topic: %s, payload: %s\n", topic, string(eventBytes))
	return nil
}

func (p *MockProducer) Close() error {
	fmt.Println("[MOCK] Producer closed")
	return nil
}
