package me

import (
	"context"
	"errors"
	"time"

	domain "github.com/umekikazuya/me/internal/domain/me"
	"github.com/umekikazuya/me/pkg/errs"
)

var _ Interactor = (*interactor)(nil)

type Interactor interface {
	UpdateProfile(ctx context.Context, in InputUpdateProfile) (*OutputDto, error)
	UpdateLinks(ctx context.Context, in InputUpdateLinks) (*OutputDto, error)
	UpdateLikes(ctx context.Context, in InputUpdateLikes) (*OutputDto, error)
	AddSkill(ctx context.Context, in InputAddSkill) (*OutputDto, error)
	RemoveSkill(ctx context.Context, in InputRemoveSkill) (*OutputDto, error)
	Get(ctx context.Context) (*OutputDto, error)
}

type interactor struct {
	repo domain.Repo
}

// AddSkill implements [Interactor].
func (i *interactor) AddSkill(ctx context.Context, in InputAddSkill) (*OutputDto, error) {
	e, err := i.repo.Find(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
		}
		return nil, errs.WrapInternal("システムエラー", err)
	}
	err = e.AddSkill(in.Name, in.Parent, time.Now())
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	err = i.repo.Save(ctx, e)
	if err != nil {
		return nil, errs.WrapInternal("システムエラー", err)
	}
	return toOutputDto(*e), nil
}

// RemoveSkill implements [Interactor].
func (i *interactor) RemoveSkill(ctx context.Context, in InputRemoveSkill) (*OutputDto, error) {
	e, err := i.repo.Find(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
		}
		return nil, errs.WrapInternal("システムエラー", err)
	}
	err = e.RemoveSkill(in.Name, in.Parent, time.Now())
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	err = i.repo.Save(ctx, e)
	if err != nil {
		return nil, errs.WrapInternal("システムエラー", err)
	}
	return toOutputDto(*e), nil
}

// UpdateLikes implements [Interactor].
func (i *interactor) UpdateLikes(ctx context.Context, in InputUpdateLikes) (*OutputDto, error) {
	e, err := i.repo.Find(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
		}
		return nil, errs.WrapInternal("システムエラー", err)
	}
	if e == nil {
		return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
	}

	err = e.UpdateLikes(in, time.Now())
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
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
			return nil, errs.New(errs.ErrBadRequest, err.Error())
		}
		links = append(links, link)
	}

	e, err := i.repo.Find(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
		}
		return nil, errs.WrapInternal("システムエラー", err)
	}
	if e == nil {
		return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
	}

	err = e.UpdateLinks(links, time.Now())
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	err = i.repo.Save(ctx, e)
	if err != nil {
		return nil, errs.WrapInternal("システムエラー", err)
	}
	return toOutputDto(*e), nil
}

// UpdateProfile implements [Interactor].
func (i *interactor) UpdateProfile(ctx context.Context, in InputUpdateProfile) (*OutputDto, error) {
	e, err := i.repo.Find(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
		}
		return nil, errs.WrapInternal("システムエラー", err)
	}
	if e == nil {
		return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
	}

	opts := make([]domain.OptProfileFunc, 0, 4)
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
) Interactor {
	return &interactor{
		repo: repo,
	}
}

func (i *interactor) Get(ctx context.Context) (*OutputDto, error) {
	e, err := i.repo.Find(ctx)
	if err != nil {
		return nil, errs.WrapInternal("me.repo.Find", err)
	}
	if e == nil {
		return nil, errs.New(errs.ErrNotFound, "Meデータが存在しません")
	}
	return toOutputDto(*e), nil
}
