package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	pg "github.com/devfullcycle/imersao22/go-gateway/internal/repository/postgres"
)

// AccountService implements domain.AccountRepository by delegating to a Postgres repository
// and also provides DTO-based methods for the API/handlers layer.
type AccountService struct {
	repo domain.AccountRepository
}

func NewAccountService(db *sql.DB) *AccountService {
	return &AccountService{repo: pg.NewPostgresAccountRepository(db)}
}

// Create creates a new account from input DTO and returns an output DTO.
func (s *AccountService) Create(ctx context.Context, in AccountCreateInput) (*AccountOutput, error) {
	acc, err := domain.NewAccount(in.Name, in.Email)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, acc); err != nil {
		return nil, err
	}
	return toAccountOutput(acc), nil
}

func (s *AccountService) GetByID(ctx context.Context, id string) (*AccountOutput, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toAccountOutput(acc), nil
}

func (s *AccountService) GetByAPIKey(ctx context.Context, apiKey string) (*AccountOutput, error) {
	fmt.Println("API Key:", apiKey)

	acc, err := s.repo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	return toAccountOutput(acc), nil
}

func (s *AccountService) UpdateBalance(ctx context.Context, apiKey string, amount float64) error {
	// Get the account by API key
	account, err := s.repo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return err
	}

	// Apply business logic: add amount to balance
	if err := account.AddBalance(amount); err != nil {
		return err
	}

	// Save the updated account to the database
	return s.repo.UpdateBalance(ctx, account)
}

// toAccountOutput maps domain.Account to output DTO.
func toAccountOutput(a *domain.Account) *AccountOutput {
	return &AccountOutput{
		ID:        a.ID,
		Name:      a.Name,
		Email:     a.Email,
		APIKey:    a.APIKey,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
