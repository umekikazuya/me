package event

import (
	"context"

	appevent "github.com/umekikazuya/me/internal/app/event"
	pkgdomain "github.com/umekikazuya/me/pkg/domain"
)

var _ appevent.EventDispatcher = (*NoopEventDispatcher)(nil)

type NoopEventDispatcher struct{}

func (n *NoopEventDispatcher) Register(_ appevent.EventHandler) {}

func (n *NoopEventDispatcher) Dispatch(_ context.Context, _ []pkgdomain.DomainEvent) error {
	return nil
}
