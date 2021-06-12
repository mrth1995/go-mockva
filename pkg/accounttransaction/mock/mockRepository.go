package mock

import (
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/model"
	"github.com/stretchr/testify/mock"
)

type MockAccountTransactionRepository struct {
	mock.Mock
}

func (r *MockAccountTransactionRepository) Save(newAccount *model.AccountTransaction) error {
	args := r.Called(newAccount)
	return args.Error(0)
}

func (r *MockAccountTransactionRepository) FindById(trxId string) (*[]model.AccountTransaction, error) {
	args := r.Called(trxId)
	var firstOutput *[]model.AccountTransaction
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*[]model.AccountTransaction)
	}
	return firstOutput, args.Error(1)
}
