package eventhandler

import (
	"context"

	appevent "github.com/umekikazuya/me/internal/app/event"
	appme "github.com/umekikazuya/me/internal/app/me"
	identitydomain "github.com/umekikazuya/me/internal/domain/identity"
	pkgdomain "github.com/umekikazuya/me/pkg/domain"
)

var _ appevent.EventHandler = (*IdentityRegisteredHandler)(nil)

type IdentityRegisteredHandler struct {
	meInteractor appme.Interactor
}

func NewIdentityRegisteredHandler(meInteractor appme.Interactor) *IdentityRegisteredHandler {
	return &IdentityRegisteredHandler{meInteractor: meInteractor}
}

func (h *IdentityRegisteredHandler) EventType() string {
	return "identity.registered"
}

func (h *IdentityRegisteredHandler) Handle(ctx context.Context, event pkgdomain.DomainEvent) error {
	e, ok := event.(identitydomain.RegisteredEvent)
	if !ok {
		return nil
	}
	_, err := h.meInteractor.Create(ctx, appme.InputDto{
		DisplayName: e.Email(),
	})
	return err
}
