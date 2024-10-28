package mock

import (
	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountTransactionRepository struct {
	mock.Mock
}

func (r *MockAccountTransactionRepository) Save(newAccount *domain.AccountTransaction) error {
	args := r.Called(newAccount)
	return args.Error(0)
}

func (r *MockAccountTransactionRepository) FindById(trxId string) (*[]domain.AccountTransaction, error) {
	args := r.Called(trxId)
	var firstOutput *[]domain.AccountTransaction
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*[]domain.AccountTransaction)
	}
	return firstOutput, args.Error(1)
}
