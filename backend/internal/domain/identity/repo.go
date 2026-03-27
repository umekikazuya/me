package identity

import "context"

type IdentityRepo interface {
	FindByID(ctx context.Context, identityID identityID) (*identityID, error)
	FindByEmail(ctx context.Context, email email) (*identityID, error)
	Save(ctx context.Context, identity *Identity) error
}

type SessionRepo interface {
	FindByUserIdAndTokenHash(
		ctx context.Context, userID string, tokenHash string,
	) (*Session, error)
	FindActiveByUser(ctx context.Context, identityID identityID) ([]*Session, error)
	Save(ctx context.Context, session *Session) error
	RevokeAll(ctx context.Context, identityID identityID) error
}