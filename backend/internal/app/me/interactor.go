package me

import (
	"context"

	domain "github.com/umekikazuya/me/internal/domain/me"
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

func (i *Interactor) Create(ctx context.Context, input InputDto) (*OutputDto, error) {
	opts := []domain.OptFunc{}
	if input.DisplayJa != nil {
		opts = append(opts, domain.OptDisplayNameJa(*input.DisplayJa))
	}
	e, err := domain.NewMe(
		input.DisplayName,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	err = i.repo.Save(e)
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
	e, err := domain.NewMe(
		input.DisplayName,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	err = i.repo.Save(e)
	if err != nil {
		return nil, err
	}

	return toOutputDto(*e), nil
}

func (i *Interactor) Get(ctx context.Context) (*OutputDto, error) {
	e, err := i.repo.Find()
	if err != nil {
		return nil, err
	}
	return toOutputDto(*e), nil
}
