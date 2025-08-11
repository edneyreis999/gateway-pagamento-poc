package domain

import "testing"

func TestNewAccount(t *testing.T) {
	a, err := NewAccount("John Doe", "john@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.ID == "" || a.APIKey == "" {
		t.Fatalf("expected IDs to be generated")
	}
	if a.Balance != 0 {
		t.Fatalf("expected initial balance 0 got %v", a.Balance)
	}
}

func TestAddBalance(t *testing.T) {
	a, _ := NewAccount("Jane", "jane@example.com")
	if err := a.AddBalance(100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Balance != 100 {
		t.Fatalf("expected balance 100 got %v", a.Balance)
	}
	if err := a.AddBalance(-10); err == nil {
		t.Fatalf("expected error for negative amount")
	}
}
