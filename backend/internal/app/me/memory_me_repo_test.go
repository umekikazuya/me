package me

import (
	"context"
	"sync"
	"testing"

	domain "github.com/umekikazuya/me/internal/domain/me"
	"github.com/umekikazuya/me/pkg/errs"
)

type memoryMeRepo struct {
	mu     sync.RWMutex
	entity *domain.Me
}

// FindByID implements [me.Repo].
func (m *memoryMeRepo) Find(ctx context.Context) (*domain.Me, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.entity == nil {
		return nil, errs.ErrNotFound
	}

	return m.entity, nil
}

// Save implements [me.Repo].
func (m *memoryMeRepo) Save(ctx context.Context, me *domain.Me) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entity = me
	return nil
}

func (m *memoryMeRepo) seedData(t *testing.T, in domain.ReconstructInput) {
	t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entity = domain.Reconstruct(in)
}

func newMeRepo() *memoryMeRepo {
	return &memoryMeRepo{}
}

var _ domain.Repo = (*memoryMeRepo)(nil)
