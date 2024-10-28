package mock

import (
	"context"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (r *MockAccountRepository) FindByID(ctx context.Context, accountId string) (*domain.Account, error) {
	args := r.Called(ctx, accountId)
	var firstOutput *domain.Account
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*domain.Account)
	}
	return firstOutput, args.Error(1)
}

func (r *MockAccountRepository) Save(ctx context.Context, newAccount *domain.Account) error {
	args := r.Called(ctx, newAccount)
	return args.Error(0)
}

func (r *MockAccountRepository) Update(ctx context.Context, updatedAccount *domain.Account) (*domain.Account, error) {
	args := r.Called(ctx, updatedAccount)
	var firstOutput *domain.Account
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*domain.Account)
	}
	return firstOutput, args.Error(1)
}

func (r *MockAccountRepository) FindAndLockAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error) {
	args := r.Called(ctx, accountID)
	var firstOutput *domain.AccountBalance
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*domain.AccountBalance)
	}
	return firstOutput, args.Error(1)
}

func (r *MockAccountRepository) UpdateBalance(ctx context.Context, accountBalance *domain.AccountBalance) (*domain.AccountBalance, error) {
	args := r.Called(ctx, accountBalance)
	var firstOutput *domain.AccountBalance
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*domain.AccountBalance)
	}
	return firstOutput, args.Error(1)
}
