package me

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/umekikazuya/me/internal/domain/me"
	"github.com/umekikazuya/me/pkg/errs"
)

var _ Interactor = (*interactor)(nil)

type Interactor interface {
	Create(ctx context.Context, input InputDto) (*OutputDto, error)
	UpdateProfile(ctx context.Context, in InputUpdateProfile) (*OutputDto, error)
	UpdateLinks(ctx context.Context, in InputUpdateLinks) (*OutputDto, error)
	UpdateLikes(ctx context.Context, in InputUpdateLikes) (*OutputDto, error)
	Get(ctx context.Context, id string) (*OutputDto, error)
}

type interactor struct {
	repo domain.Repo
	id   string
}

// UpdateLikes implements [Interactor].
func (i *interactor) UpdateLikes(ctx context.Context, in InputUpdateLikes) (*OutputDto, error) {
	e, err := i.repo.FindByID(ctx, i.id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
		}
		return nil, errs.WrapInternal("システムエラー", err)
	}

	err = e.UpdateLikes(in, time.Now())
	if err != nil {
		return nil, errs.New(errs.ErrConflict, err.Error())
	}
	err = i.repo.Save(ctx, e)
	if err != nil {
		return nil, errs.WrapInternal("システムエラー", err)
	}
	return toOutputDto(*e), nil
}

// UpdateLinks implements [Interactor].
func (i *interactor) UpdateLinks(ctx context.Context, in InputUpdateLinks) (*OutputDto, error) {
	links := make([]domain.Link, 0, len(in))
	for _, l := range in {
		link, err := domain.NewLink(l.Platform, l.URL)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	e, err := i.repo.FindByID(ctx, i.id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
		}
		return nil, errs.WrapInternal("システムエラー", err)
	}

	err = e.UpdateLinks(links, time.Now())
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	return toOutputDto(*e), nil
}

// UpdateProfile implements [Interactor].
func (i *interactor) UpdateProfile(ctx context.Context, in InputUpdateProfile) (*OutputDto, error) {
	e, err := i.repo.FindByID(ctx, i.id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
		}
		return nil, errs.WrapInternal("システムエラー", err)
	}
	opts := make([]domain.OptProfileFunc, 0)
	opts = append(opts, domain.OptDisplayName(in.DisplayName))
	opts = append(opts, domain.OptDisplayNameJa(in.DisplayJa))
	opts = append(opts, domain.OptRole(in.Role))
	opts = append(opts, domain.OptLocation(in.Location))

	err = e.UpdateProfile(time.Now(), opts...)
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	err = i.repo.Save(ctx, e)
	if err != nil {
		return nil, errs.WrapInternal("システムエラー", err)
	}
	return toOutputDto(*e), nil
}

// NewInteractor はユースケースの初期化クラス
func NewInteractor(
	repo domain.Repo,
	id string,
) Interactor {
	return &interactor{
		repo: repo,
		id:   id,
	}
}

func (i *interactor) Create(ctx context.Context, input InputDto) (*OutputDto, error) {
	exists, err := i.repo.Exists(ctx, input.ID)
	if err != nil {
		return nil, errs.WrapInternal("me.repo.Exists", err)
	}
	if exists {
		return nil, fmt.Errorf("create me: %w", errs.ErrConflict)
	}

	e, err := domain.NewMe(
		input.ID,
	)
	if err != nil {
		return nil, err
	}

	err = i.repo.Save(ctx, e)
	if err != nil {
		return nil, errs.WrapInternal("me.repo.Save", err)
	}

	return toOutputDto(*e), nil
}

func (i *interactor) Get(ctx context.Context, id string) (*OutputDto, error) {
	e, err := i.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.WrapInternal("me.repo.FindByID", err)
	}
	if e == nil {
		return nil, fmt.Errorf("get me: %w", errs.ErrNotFound)
	}
	return toOutputDto(*e), nil
}
