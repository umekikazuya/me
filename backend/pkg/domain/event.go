package domain

import "time"

type DomainEvent interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
}
