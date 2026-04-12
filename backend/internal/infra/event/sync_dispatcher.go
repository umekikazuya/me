package event

import (
	"context"

	appevent "github.com/umekikazuya/me/internal/app/event"
	pkgdomain "github.com/umekikazuya/me/pkg/domain"
)

var _ appevent.EventDispatcher = (*SyncEventDispatcher)(nil)

type SyncEventDispatcher struct {
	handlers map[string][]appevent.EventHandler
}

func NewSyncEventDispatcher() *SyncEventDispatcher {
	return &SyncEventDispatcher{
		handlers: make(map[string][]appevent.EventHandler),
	}
}

func (d *SyncEventDispatcher) Register(h appevent.EventHandler) {
	d.handlers[h.EventType()] = append(d.handlers[h.EventType()], h)
}

func (d *SyncEventDispatcher) Dispatch(ctx context.Context, events []pkgdomain.DomainEvent) error {
	for _, e := range events {
		for _, h := range d.handlers[e.EventType()] {
			if err := h.Handle(ctx, e); err != nil {
				return err
			}
		}
	}
	return nil
}
