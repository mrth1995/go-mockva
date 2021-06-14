package account

import (
	"github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type Repository interface {
	FindById(accountId string) (*Account, error)
	Save(newAccount *Account) error
	Update(updatedAccount *Account) (*Account, error)
}

type repositoryImpl struct {
	Connection *gorm.DB
}

func (r *repositoryImpl) FindById(accountId string) (*Account, error) {
	var existingUser Account
	find := r.Connection.First(&existingUser, "id = ?", accountId)
	if find.Error != nil && find.Error == gorm.ErrRecordNotFound {
		return nil, errors.NewAccountNotFound(accountId)
	}
	if find.Error != nil {
		return nil, find.Error
	}
	return &existingUser, nil
}

func (r *repositoryImpl) Save(newAccount *Account) error {
	tx := r.Connection.Begin()
	tx.Create(newAccount)
	tx.Commit()
	return nil
}

func (r *repositoryImpl) Update(updatedAccount *Account) (*Account, error) {
	tx := r.Connection.Begin()
	tx.Save(updatedAccount)
	tx.Commit()
	return updatedAccount, nil
}

func NewRepository(dbConnection *gorm.DB) Repository {
	return &repositoryImpl{
		Connection: dbConnection,
	}
}

type MockRepository struct {
	mock.Mock
}

func (r *MockRepository) FindById(accountId string) (*Account, error) {
	args := r.Called(accountId)
	var firstOutput *Account
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*Account)
	}
	return firstOutput, args.Error(1)
}

func (r *MockRepository) Save(newAccount *Account) error {
	args := r.Called(newAccount)
	return args.Error(0)
}

func (r *MockRepository) Update(updatedAccount *Account) (*Account, error) {
	args := r.Called(updatedAccount)
	var firstOutput *Account
	if args.Get(0) != nil {
		firstOutput = args.Get(0).(*Account)
	}
	return firstOutput, args.Error(1)
}
