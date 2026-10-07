package identity

import (
	"time"

	"github.com/google/uuid"
)

// Account は認証・認可の集約
type (
	Account struct {
		id       id
		provider struct {
			github *github
		}
		createdAt time.Time
		updatedAt time.Time
	}
	github struct {
		id string
	}
)

// RegisterWithGithub returns Account profile
//
// github id を受け取ってプロバイダにセット
func RegisterWithGithub(id string) (*Account, error) {
	return nil, nil
}

// --- Reconstruct ---

// ReconstructIdentityInput はReconstructIdentityの入力型
type ReconstructIdentityInput struct {
	ID        uuid.UUID
	githubID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ReconstructIdentity はDBから取得した信頼済みデータでIdentityを復元する
func ReconstructIdentity(input ReconstructIdentityInput) (*Account, error) {
	e := &Account{
		id:        id{value: input.ID},
		createdAt: input.CreatedAt,
		updatedAt: input.UpdatedAt,
	}
	if input.githubID != "" {
		e.provider.github.id = input.githubID
	}
	return e, nil
}

// --- Getter ---

// ID returns unique id
func (e *Account) ID() string {
	return e.id.Value()
}

// CreatedAt returnes when entity was created
func (e *Account) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt returnes when entity was updated
func (e *Account) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Account) GithubID() string {
	return e.provider.github.id
}
