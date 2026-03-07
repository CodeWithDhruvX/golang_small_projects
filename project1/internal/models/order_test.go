package models

import (
	"encoding/json"
	"testing"
)

func TestOrderJSON(t *testing.T) {
	jsonBody := `{"item":"Laptop","quantity":1,"price":1200.5}`
	var order Order
	err := json.Unmarshal([]byte(jsonBody), &order)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if order.Item != "Laptop" {
		t.Errorf("Expected item 'Laptop', got %s", order.Item)
	}
	if order.LegacyQty != 1 {
		t.Errorf("Expected quantity 1, got %d", order.LegacyQty)
	}
	if order.LegacyPrice != 1200.5 {
		t.Errorf("Expected price 1200.5, got %f", order.LegacyPrice)
	}

	// Test Marshal
	order.ID = "12345"
	order.Status = "PENDING"
	bytes, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	if len(bytes) == 0 {
		t.Fatalf("Expected non-empty JSON bytes")
	}
}
