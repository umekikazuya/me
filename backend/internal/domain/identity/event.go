package identity

import "time"

type RegisteredEvent struct {
	identityID string
	email      string
	occurredAt time.Time
}

func (e RegisteredEvent) EventType() string     { return "identity.registered" }
func (e RegisteredEvent) AggregateID() string   { return e.identityID }
func (e RegisteredEvent) OccurredAt() time.Time { return e.occurredAt }
func (e RegisteredEvent) IdentityID() string    { return e.identityID }
func (e RegisteredEvent) Email() string         { return e.email }

type AuthenticatedEvent struct {
	identityID string
	email      string
	occurredAt time.Time
}

func (e AuthenticatedEvent) EventType() string     { return "identity.authenticated" }
func (e AuthenticatedEvent) AggregateID() string   { return e.identityID }
func (e AuthenticatedEvent) OccurredAt() time.Time { return e.occurredAt }

type PasswordResetEvent struct {
	identityID string
	email      string
	occurredAt time.Time
}

func (e PasswordResetEvent) EventType() string     { return "identity.passwordReset" }
func (e PasswordResetEvent) AggregateID() string   { return e.identityID }
func (e PasswordResetEvent) OccurredAt() time.Time { return e.occurredAt }

type EmailChangedEvent struct {
	identityID string
	email      string
	occurredAt time.Time
}

func (e EmailChangedEvent) EventType() string     { return "identity.emailChanged" }
func (e EmailChangedEvent) AggregateID() string   { return e.identityID }
func (e EmailChangedEvent) OccurredAt() time.Time { return e.occurredAt }

type SessionRevokedEvent struct {
	identityID string
	occurredAt time.Time
}

func (e SessionRevokedEvent) EventType() string     { return "identity.sessionRevoked" }
func (e SessionRevokedEvent) AggregateID() string   { return e.identityID }
func (e SessionRevokedEvent) OccurredAt() time.Time { return e.occurredAt }

type SessionRotateEvent struct {
	identityID string
	occurredAt time.Time
}

func (e SessionRotateEvent) EventType() string     { return "identity.sessionRotated" }
func (e SessionRotateEvent) AggregateID() string   { return e.identityID }
func (e SessionRotateEvent) OccurredAt() time.Time { return e.occurredAt }
