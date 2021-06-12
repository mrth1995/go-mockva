package mock

import (
	"github.com/mrth1995/go-mockva/pkg/account/model"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (r *MockAccountRepository) FindById(accountId string) (*model.Account, error) {
	args := r.Called(accountId)
	var firstOutput *model.Account
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*model.Account)
	}
	return firstOutput, args.Error(1)
}

func (r *MockAccountRepository) Save(newAccount *model.Account) error {
	args := r.Called(newAccount)
	return args.Error(0)
}

func (r *MockAccountRepository) Update(updatedAccount *model.Account) (*model.Account, error) {
	args := r.Called(updatedAccount)
	var firstOutput *model.Account
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*model.Account)
	}
	return firstOutput, args.Error(1)
}
