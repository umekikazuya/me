package identity

type (
	EventType   string
	DomainEvent interface {
		Type() EventType
	}
)

const (
	EventTypeRegistered     EventType = "registered"
	EventTypeAuthenticated  EventType = "authenticated"
	EventTypePasswordReset  EventType = "passwordReset"
	EventTypeEmailChanged   EventType = "emailChanged"
	EventTypeSessionRevoked EventType = "sessionRevoked"
	EventTypeSessionRotated EventType = "sessionRotated"
)

func (e EventType) Type() EventType {
	return e
}
