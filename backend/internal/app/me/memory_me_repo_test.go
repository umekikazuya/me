package me

import (
	"context"
	"sync"
	"testing"

	domain "github.com/umekikazuya/me/internal/domain/me"
	"github.com/umekikazuya/me/pkg/errs"
)

type memoryMeRepo struct {
	mu sync.RWMutex
	es map[string]*domain.Me
}

// Exists implements [me.Repo].
func (m *memoryMeRepo) Exists(ctx context.Context, id string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exist := m.es[id]
	return exist, nil
}

// FindByID implements [me.Repo].
func (m *memoryMeRepo) FindByID(ctx context.Context, id string) (*domain.Me, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	e, exist := m.es[id]
	if !exist {
		return nil, errs.ErrNotFound
	}
	return e, nil
}

// Save implements [me.Repo].
func (m *memoryMeRepo) Save(ctx context.Context, me *domain.Me) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.es[me.ID()] = me
	return nil
}

func (m *memoryMeRepo) seedData(t *testing.T, in domain.ReconstructInput) {
	t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()

	m.es[in.ID.String()] = domain.Reconstruct(in)
}

func newMeRepo() *memoryMeRepo {
	return &memoryMeRepo{
		es: make(map[string]*domain.Me),
	}
}

var _ domain.Repo = (*memoryMeRepo)(nil)
