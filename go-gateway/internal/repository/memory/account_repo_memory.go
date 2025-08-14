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

func (r *InMemoryAccountRepository) UpdateBalance(ctx context.Context, a *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Find the account in memory
	storedAccount, ok := r.byID[a.ID]
	if !ok {
		return domain.ErrAccountNotFound
	}

	// Update the stored account with the new balance and updated_at
	storedAccount.Balance = a.Balance
	storedAccount.UpdatedAt = a.UpdatedAt

	// Also update the API key map
	if storedAccount.APIKey != a.APIKey {
		delete(r.byAPIKey, storedAccount.APIKey)
		r.byAPIKey[a.APIKey] = storedAccount
	}

	return nil
}
