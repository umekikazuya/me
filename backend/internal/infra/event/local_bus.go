package event

import (
	"context"

	appevent "github.com/umekikazuya/me/internal/app/event"
	pkgdomain "github.com/umekikazuya/me/pkg/domain"
)

var _ appevent.EventDispatcher = (*LocalEventBus)(nil)

type LocalEventBus struct {
	handlers map[string][]appevent.EventHandler
}

func NewLocalEventBus() *LocalEventBus {
	return &LocalEventBus{
		handlers: make(map[string][]appevent.EventHandler),
	}
}

func (b *LocalEventBus) Register(h appevent.EventHandler) {
	b.handlers[h.EventType()] = append(b.handlers[h.EventType()], h)
}

func (b *LocalEventBus) Dispatch(ctx context.Context, events []pkgdomain.DomainEvent) error {
	for _, e := range events {
		for _, h := range b.handlers[e.EventType()] {
			if err := h.Handle(ctx, e); err != nil {
				return err
			}
		}
	}
	return nil
}
