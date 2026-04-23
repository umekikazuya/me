package me

import (
	"context"
	"fmt"

	domain "github.com/umekikazuya/me/internal/domain/me"
	"github.com/umekikazuya/me/pkg/errs"
)

var _ Interactor = (*interactor)(nil)

type Interactor interface {
	Create(ctx context.Context, input InputDto) (*OutputDto, error)
	Update(ctx context.Context, input InputDto) (*OutputDto, error)
	Get(ctx context.Context, id string) (*OutputDto, error)
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

func (i *interactor) Create(ctx context.Context, input InputDto) (*OutputDto, error) {
	exists, err := i.repo.Exists(ctx, input.ID)
	if err != nil {
		return nil, errs.WrapInternal(ctx, "me.repo.Exists", err)
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
	if input.Links != nil {
		links := make([]domain.Link, 0, len(input.Links))
		for _, l := range input.Links {
			link, err := domain.NewLink(l.Platform, l.URL)
			if err != nil {
				return nil, err
			}
			links = append(links, link)
		}
		opts = append(opts, domain.OptLinks(links))
	}
	if input.Certifications != nil {
		certs := make([]domain.Certification, 0, len(input.Certifications))
		for _, c := range input.Certifications {
			cert, err := domain.NewCertification(c.Name, c.Issuer, c.Year, c.Month)
			if err != nil {
				return nil, err
			}
			certs = append(certs, cert)
		}
		opts = append(opts, domain.OptCertifications(certs))
	}
	e, err := domain.NewMe(
		input.ID,
		input.DisplayName,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	err = i.repo.Save(ctx, e)
	if err != nil {
		return nil, errs.WrapInternal(ctx, "me.repo.Save", err)
	}

	return toOutputDto(*e), nil
}

func (i *interactor) Update(ctx context.Context, input InputDto) (*OutputDto, error) {
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
	if input.Links != nil {
		links := make([]domain.Link, 0, len(input.Links))
		for _, l := range input.Links {
			link, err := domain.NewLink(l.Platform, l.URL)
			if err != nil {
				return nil, err
			}
			links = append(links, link)
		}
		opts = append(opts, domain.OptLinks(links))
	}
	if input.Certifications != nil {
		certs := make([]domain.Certification, 0, len(input.Certifications))
		for _, c := range input.Certifications {
			cert, err := domain.NewCertification(c.Name, c.Issuer, c.Year, c.Month)
			if err != nil {
				return nil, err
			}
			certs = append(certs, cert)
		}
		opts = append(opts, domain.OptCertifications(certs))
	}
	e, err := i.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, errs.WrapInternal(ctx, "me.repo.FindByID", err)
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
		return nil, errs.WrapInternal(ctx, "me.repo.Save", err)
	}

	return toOutputDto(*e), nil
}

func (i *interactor) Get(ctx context.Context, id string) (*OutputDto, error) {
	e, err := i.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.WrapInternal(ctx, "me.repo.FindByID", err)
	}
	if e == nil {
		return nil, fmt.Errorf("get me: %w", errs.ErrNotFound)
	}
	return toOutputDto(*e), nil
}
