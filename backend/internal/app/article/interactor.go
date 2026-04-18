package article

import (
	"context"
	"errors"

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
	Sync(ctx context.Context, platform string) domain.IndexingResult
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

func (i *interactor) Search(ctx context.Context, input InputSearchDto) (*OutputSearchDto, error) {
	return nil, errors.New("not implemented")
}

func (i *interactor) GetTagsAll(ctx context.Context) (*OutputTagAllDto, error) {
	return nil, errors.New("not implemented")
}

func (i *interactor) GetSuggests(ctx context.Context, input InputGetSuggestDto) (*OutputGetSuggestAllDto, error) {
	return nil, errors.New("not implemented")
}

func (i *interactor) Register(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (i *interactor) Update(ctx context.Context, input InputUpdateDto) error {
	return errors.New("not implemented")
}

func (i *interactor) Remove(ctx context.Context, input InputRemoveDto) error {
	return errors.New("not implemented")
}

func (i *interactor) Sync(ctx context.Context, platform string) domain.IndexingResult {
	return domain.IndexingResult{}
}
