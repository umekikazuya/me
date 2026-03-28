package identity

import "context"

type IdentityRepo interface {
	FindByID(ctx context.Context, id string) (*Identity, error)
	FindByEmail(ctx context.Context, email email) (*Identity, error)
	Save(ctx context.Context, identity *Identity) error
}

type SessionRepo interface {
	FindByIdentityIdAndTokenHash(
		ctx context.Context, identityID string, tokenHash string,
	) (*Session, error)
	FindActiveByIdentity(ctx context.Context, identityID identityID) ([]*Session, error)
	Save(ctx context.Context, session *Session) error
	RevokeAll(ctx context.Context, id string) error
}
