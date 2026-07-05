package domain

import "encoding/json"

type Order struct {
	ID string
	CustomerID string
	ProductID string
	Quantity int
	Amount float64
}

type OutboxEvent struct {
    ID            string
    AggregateType string
    AggregateID   string
    EventType     string
    Payload       json.RawMessage
    Status        string
}