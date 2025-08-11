package memory

import (
	"context"
	"sync"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

// InMemoryAccountRepository is a thread-safe in-memory repository.
type InMemoryAccountRepository struct {
	mu       sync.RWMutex
	byID     map[string]*domain.Account
	byAPIKey map[string]*domain.Account
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		byID:     make(map[string]*domain.Account),
		byAPIKey: make(map[string]*domain.Account),
	}
}

func (r *InMemoryAccountRepository) Create(ctx context.Context, a *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[a.ID] = a
	r.byAPIKey[a.APIKey] = a
	return nil
}

func (r *InMemoryAccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if a, ok := r.byID[id]; ok {
		return a, nil
	}
	return nil, domain.ErrAccountNotFound
}

func (r *InMemoryAccountRepository) GetByAPIKey(ctx context.Context, apiKey string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if a, ok := r.byAPIKey[apiKey]; ok {
		return a, nil
	}
	return nil, domain.ErrAccountNotFound
}

func (r *InMemoryAccountRepository) UpdateBalance(ctx context.Context, id string, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.byID[id]
	if !ok {
		return domain.ErrAccountNotFound
	}
	return a.AddBalance(amount)
}
