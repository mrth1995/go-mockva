package repository

import (
	"context"

	"github.com/mrth1995/go-mockva/pkg/domain"
)

type AccountRepository interface {
	FindByID(ctx context.Context, accountId string) (*domain.Account, error)
	Save(ctx context.Context, newAccount *domain.Account) error
	Update(ctx context.Context, updatedAccount *domain.Account) (*domain.Account, error)
	FindAndLockAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error)
	UpdateBalance(ctx context.Context, accountBalance *domain.AccountBalance) (*domain.AccountBalance, error)
}
