package repository

//go:generate mockgen -destination=mock/mockAccountRepository.go -package=mock github.com/mrth1995/go-mockva/pkg/repository AccountRepository

import (
	"context"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"gorm.io/gorm"
)

// AccountRepository defines the interface for account data persistence operations.
type AccountRepository interface {
	// FindByID retrieves an account by its unique identifier.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - accountId: The unique account identifier
	// Returns:
	//   - *domain.Account: The account if found
	//   - error: If the account is not found or a database error occurs
	FindByID(ctx context.Context, accountId string) (*domain.Account, error)

	// Save persists a new account to the database.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - newAccount: The account entity to persist
	// Returns:
	//   - error: If the account already exists or a database error occurs
	Save(ctx context.Context, newAccount *domain.Account) error

	// Update modifies an existing account's information in the database.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - updatedAccount: The account entity with updated fields
	// Returns:
	//   - *domain.Account: The updated account
	//   - error: If the account is not found or a database error occurs
	Update(ctx context.Context, updatedAccount *domain.Account) (*domain.Account, error)

	// FindAndLockAccountBalance retrieves an account balance with a pessimistic lock.
	// This method should be used within transactions to prevent concurrent balance modifications.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - accountID: The unique account identifier
	// Returns:
	//   - *domain.AccountBalance: The account balance with an active row lock
	//   - error: If the account is not found or a database error occurs
	FindAndLockAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error)

	// UpdateBalance updates the account balance within the provided transaction context.
	// This method must be called within an active database transaction.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - accountBalance: The AccountBalance entity with the new balance value
	//   - tx: The GORM transaction context
	// Returns:
	//   - *domain.AccountBalance: The updated balance
	//   - error: If the account is not found or a database error occurs
	UpdateBalance(ctx context.Context, accountBalance *domain.AccountBalance, tx *gorm.DB) (*domain.AccountBalance, error)
}
