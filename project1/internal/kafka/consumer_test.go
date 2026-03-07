package kafka

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

// mockConsumerGroupSession is a mock implementation of sarama.ConsumerGroupSession
type mockConsumerGroupSession struct {
	ctx context.Context
}

func (m *mockConsumerGroupSession) Claims() map[string][]int32 { return nil }
func (m *mockConsumerGroupSession) MemberID() string           { return "test" }
func (m *mockConsumerGroupSession) GenerationID() int32        { return 1 }
func (m *mockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}
func (m *mockConsumerGroupSession) Commit() {}
func (m *mockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}
func (m *mockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {}
func (m *mockConsumerGroupSession) Context() context.Context                                 { return m.ctx }

// mockConsumerGroupClaim is a mock implementation of sarama.ConsumerGroupClaim
type mockConsumerGroupClaim struct {
	messages chan *sarama.ConsumerMessage
}

func (m *mockConsumerGroupClaim) Topic() string                            { return "orders.created" }
func (m *mockConsumerGroupClaim) Partition() int32                         { return 0 }
func (m *mockConsumerGroupClaim) InitialOffset() int64                     { return 0 }
func (m *mockConsumerGroupClaim) HighWaterMarkOffset() int64               { return 0 }
func (m *mockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage { return m.messages }

func TestConsumerGroupHandler_ConsumeClaim(t *testing.T) {
	processed := false
	var wg sync.WaitGroup
	wg.Add(1)

	processFunc := func(msg *sarama.ConsumerMessage) error {
		defer wg.Done()
		if string(msg.Value) != "test message" {
			t.Errorf("Expected 'test message', got '%s'", string(msg.Value))
		}
		processed = true
		return nil
	}

	handler := &ConsumerGroupHandler{
		Ready:   make(chan bool),
		Process: processFunc,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := &mockConsumerGroupSession{ctx: ctx}

	msgChan := make(chan *sarama.ConsumerMessage, 1)
	msgChan <- &sarama.ConsumerMessage{Value: []byte("test message")}

	claim := &mockConsumerGroupClaim{messages: msgChan}

	go func() {
		err := handler.ConsumeClaim(session, claim)
		if err != nil {
			t.Errorf("ConsumeClaim returned error: %v", err)
		}
	}()

	// Wait for process func to finish or timeout
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		// Done
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for message to process")
	}

	if !processed {
		t.Error("Message was not processed")
	}
}
