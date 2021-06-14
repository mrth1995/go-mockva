package accounttransaction

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type Repository interface {
	Save(trx *AccountTransaction) error
}

type repositoryImpl struct {
	Connection *gorm.DB
}

func (r *repositoryImpl) Save(trx *AccountTransaction) error {
	tx := r.Connection.Begin()
	tx.Create(trx)
	tx.Commit()
	return nil
}

func NewRepository(dbConnection *gorm.DB) Repository {
	return &repositoryImpl{
		Connection: dbConnection,
	}
}

type MockRepository struct {
	mock.Mock
}

func (r *MockRepository) Save(newAccount *AccountTransaction) error {
	args := r.Called(newAccount)
	return args.Error(0)
}

func (r *MockRepository) FindById(trxId string) (*[]AccountTransaction, error) {
	args := r.Called(trxId)
	var firstOutput *[]AccountTransaction
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*[]AccountTransaction)
	}
	return firstOutput, args.Error(1)
}
