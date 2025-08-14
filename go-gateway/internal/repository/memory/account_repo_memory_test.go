package memory

import (
	"context"
	"testing"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

func TestInMemoryAccountRepository(t *testing.T) {
	repo := NewInMemoryAccountRepository()
	ctx := context.Background()

	a, err := domain.NewAccount("Acme", "acme@example.com")
	if err != nil {
		t.Fatalf("new account: %v", err)
	}

	if err := repo.Create(ctx, a); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.GetByID(ctx, a.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.ID != a.ID {
		t.Fatalf("expected same id")
	}

	got2, err := repo.GetByAPIKey(ctx, a.APIKey)
	if err != nil {
		t.Fatalf("get by apikey: %v", err)
	}
	if got2.APIKey != a.APIKey {
		t.Fatalf("expected same apikey")
	}

	// Test UpdateBalance with account
	updatedAccount := &domain.Account{
		ID:        a.ID,
		Name:      a.Name,
		Email:     a.Email,
		APIKey:    a.APIKey,
		Balance:   50.0,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	if err := repo.UpdateBalance(ctx, updatedAccount); err != nil {
		t.Fatalf("update balance: %v", err)
	}

	// Verify the balance was updated
	got3, err := repo.GetByID(ctx, a.ID)
	if err != nil {
		t.Fatalf("get by id after update: %v", err)
	}
	if got3.Balance != 50.0 {
		t.Fatalf("expected balance 50, got %v", got3.Balance)
	}

	if _, err := repo.GetByID(ctx, "does-not-exist"); err == nil {
		t.Fatalf("expected not found error")
	}
}
