package article

import (
	"context"

	domain "github.com/umekikazuya/me/internal/domain/article"
)

var _ Interactor = (*interactor)(nil)

type Interactor interface {
	Search(ctx context.Context, input InputSearchDto) (*OutputSearchDto, error)
	GetTagsAll(ctx context.Context) (*OutputTagAllDto, error)
	GetSuggests(ctx context.Context, input InputGetSuggestDto) (*OutputGetSuggestAllDto, error)
	Register(ctx context.Context, id string) error
	Update(ctx context.Context, input InputUpdateDto) error
	Remove(ctx context.Context, input InputRemoveDto) error
	Index(ctx context.Context) error
}

type interactor struct {
	repo domain.Repo
}

// NewInteractor はユースケースの初期化クラス
func NewInteractor(
	repo domain.Repo,
) Interactor {
	return &interactor{
		repo: repo,
	}
}
