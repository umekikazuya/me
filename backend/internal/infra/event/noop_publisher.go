package event

import (
	"context"

	"github.com/umekikazuya/me/pkg/domain"
)

type NoopEventPublisher struct{}

func (n *NoopEventPublisher) Publish(_ context.Context, _ []domain.DomainEvent) error {
	return nil
}
