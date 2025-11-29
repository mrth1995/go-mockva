package repository

//go:generate mockgen -destination=mock/mockAccountRepository.go -package=mock github.com/mrth1995/go-mockva/pkg/repository AccountRepository

import (
	"context"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"gorm.io/gorm"
)

type AccountRepository interface {
	FindByID(ctx context.Context, accountId string) (*domain.Account, error)
	Save(ctx context.Context, newAccount *domain.Account) error
	Update(ctx context.Context, updatedAccount *domain.Account) (*domain.Account, error)
	FindAndLockAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error)
	// UpdateBalance updates the account balance within the provided transaction context.
	// Parameters:
	//   - ctx: The request context
	//   - accountBalance: The AccountBalance to update
	//   - tx: The GORM transaction context
	// Returns:
	//   - *domain.AccountBalance: The updated balance
	//   - error: If the operation fails
	UpdateBalance(ctx context.Context, accountBalance *domain.AccountBalance, tx *gorm.DB) (*domain.AccountBalance, error)
}
