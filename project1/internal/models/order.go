package models

import (
	"encoding/json"
)

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// LegacyOrder represents the old format for backward compatibility
type LegacyOrder struct {
	Item     string  `json:"item"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	ID           string      `json:"id"`
	CustomerID   string      `json:"customer_id,omitempty"`
	Items        []OrderItem `json:"items,omitempty"`
	TotalAmount  float64     `json:"total_amount,omitempty"`
	
	// Legacy fields for backward compatibility
	Item         string      `json:"item,omitempty"`
	LegacyQty    int         `json:"quantity,omitempty"`
	LegacyPrice  float64     `json:"price,omitempty"`
	
	Status       string      `json:"status"`
	CreatedAt    string      `json:"created_at"`
}

// UnmarshalJSON handles both old and new formats
func (o *Order) UnmarshalJSON(data []byte) error {
	// Try new format first
	type NewOrder Order
	var newOrder NewOrder
	if err := json.Unmarshal(data, &newOrder); err == nil && newOrder.CustomerID != "" {
		*o = Order(newOrder)
		return nil
	}
	
	// Try legacy format
	var legacy LegacyOrder
	if err := json.Unmarshal(data, &legacy); err == nil {
		o.CustomerID = "legacy-customer"
		o.Items = []OrderItem{{
			ProductID: legacy.Item,
			Quantity:  legacy.Quantity,
			Price:     legacy.Price,
		}}
		o.TotalAmount = legacy.Price * float64(legacy.Quantity)
		o.Item = legacy.Item
		o.LegacyQty = legacy.Quantity
		o.LegacyPrice = legacy.Price
		return nil
	}
	
	return json.Unmarshal(data, (*NewOrder)(o))
}
