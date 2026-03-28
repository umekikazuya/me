package domain

import (
	"context"
	"time"
)

type DomainEvent interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
}

type EventPublisher interface {
	Publish(ctx context.Context, events []DomainEvent) error
}
