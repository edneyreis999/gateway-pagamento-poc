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

	if err := repo.UpdateBalance(ctx, a.ID, 50); err != nil {
		t.Fatalf("update balance: %v", err)
	}
	if a.Balance != 50 {
		t.Fatalf("expected balance 50, got %v", a.Balance)
	}

	if _, err := repo.GetByID(ctx, "does-not-exist"); err == nil {
		t.Fatalf("expected not found error")
	}
}
