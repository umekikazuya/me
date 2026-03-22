package me

import (
	"context"
	"fmt"

	domain "github.com/umekikazuya/me/internal/domain/me"
	"github.com/umekikazuya/me/pkg/errs"
)

var _ interactor = (*Interactor)(nil)

type interactor interface {
	Create(ctx context.Context, input InputDto) (*OutputDto, error)
	Update(ctx context.Context, input InputDto) (*OutputDto, error)
	Get(ctx context.Context) (*OutputDto, error)
}

type Interactor struct {
	repo domain.Repo
}

// NewInteractor はユースケースの初期化クラス
func NewInteractor(
	repo domain.Repo,
) interactor {
	return &Interactor{
		repo: repo,
	}
}

func (i *Interactor) Create(ctx context.Context, input InputDto) (*OutputDto, error) {
	exists, err := i.repo.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("create me: %w", errs.ErrConflict)
	}

	opts := []domain.OptFunc{}
	if input.DisplayJa != nil {
		opts = append(opts, domain.OptDisplayNameJa(*input.DisplayJa))
	}
	if input.Role != nil {
		opts = append(opts, domain.OptRole(*input.Role))
	}
	if input.Location != nil {
		opts = append(opts, domain.OptLocation(*input.Location))
	}
	if input.Likes != nil {
		opts = append(opts, domain.OptLikes(input.Likes))
	}
	e, err := domain.NewMe(
		input.DisplayName,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	err = i.repo.Save(ctx, e)
	if err != nil {
		return nil, err
	}

	return toOutputDto(*e), nil
}

func (i *Interactor) Update(ctx context.Context, input InputDto) (*OutputDto, error) {
	opts := []domain.OptFunc{}
	if input.DisplayJa != nil {
		opts = append(opts, domain.OptDisplayNameJa(*input.DisplayJa))
	}
	if input.Role != nil {
		opts = append(opts, domain.OptRole(*input.Role))
	}
	if input.Location != nil {
		opts = append(opts, domain.OptLocation(*input.Location))
	}
	if input.Likes != nil {
		opts = append(opts, domain.OptLikes(input.Likes))
	}
	e, err := i.repo.Find(ctx)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, fmt.Errorf("update me: %w", errs.ErrNotFound)
	}
	err = e.Update(input.DisplayName, opts...)
	if err != nil {
		return nil, fmt.Errorf("update me: %w: %w", errs.ErrUnprocessable, err)
	}

	err = i.repo.Save(ctx, e)
	if err != nil {
		return nil, err
	}

	return toOutputDto(*e), nil
}

func (i *Interactor) Get(ctx context.Context) (*OutputDto, error) {
	e, err := i.repo.Find(ctx)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, fmt.Errorf("get me: %w", errs.ErrNotFound)
	}
	return toOutputDto(*e), nil
}
