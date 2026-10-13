package identity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Account は認証・認可の集約
type (
	Account struct {
		id        id
		provider  provider
		createdAt time.Time
		updatedAt time.Time
	}
	provider struct {
		github *github
	}
	github struct {
		id string
	}
)

// RegisterWithGithub returns Account profile
//
// github id を受け取ってプロバイダにセット
func RegisterWithGithub(
	inputID string,
	baseTime time.Time,
) (*Account, error) {
	if inputID == "" {
		return nil, errors.New("github 認証が不正です")
	}
	return &Account{
		id: id{
			value: uuid.New(),
		},
		provider: provider{
			github: &github{
				id: inputID,
			},
		},
		createdAt: baseTime,
		updatedAt: baseTime,
	}, nil
}

// --- Reconstruct ---

// ReconstructIdentityInput はReconstructIdentityの入力型
type ReconstructIdentityInput struct {
	InputID        uuid.UUID
	InputGithubID  string
	InputCreatedAt time.Time
	InputUpdatedAt time.Time
}

// ReconstructAccount はDBから取得した信頼済みデータでIdentityを復元する
func ReconstructAccount(input ReconstructIdentityInput) (*Account, error) {
	e := &Account{
		id:        id{value: input.InputID},
		createdAt: input.InputCreatedAt,
		updatedAt: input.InputUpdatedAt,
		provider: provider{
			github: &github{},
		},
	}
	if input.InputGithubID != "" {
		e.provider.github.id = input.InputGithubID
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
