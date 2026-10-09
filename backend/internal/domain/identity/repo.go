package identity

import "context"

type IdentityRepo interface {
	FindByID(ctx context.Context, id string) (*Account, error)
	FindByEmail(ctx context.Context, email string) (*Account, error)
	Save(ctx context.Context, identity *Account) error
	FindByGithubID(ctx context.Context, githubID string) (*Account, error)
}

// TODO: セッションのローテーションを原始的にするよう RotateSessions(ctx context.Context, old, new *domain.Session) error を追加する
type SessionRepo interface {
	FindByIdentityIdAndTokenHash(
		ctx context.Context, identityID string, tokenHash string,
	) (*Session, error)
	FindActiveByIdentity(ctx context.Context, identityID string) ([]*Session, error)
	Save(ctx context.Context, session *Session) error
	RevokeAll(ctx context.Context, id string) error
}
