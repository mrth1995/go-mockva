package repository

import (
	"github.com/mrth1995/go-mockva/pkg/errors"

	"github.com/mrth1995/go-mockva/pkg/account/model"
	"gorm.io/gorm"
)

type AccountRepositoryImpl struct {
	Connection *gorm.DB
}

func (r *AccountRepositoryImpl) FindById(accountId string) (*model.Account, error) {
	var existingUser model.Account
	find := r.Connection.First(&existingUser, "id = ?", accountId)
	if find.Error != nil && find.Error == gorm.ErrRecordNotFound {
		return nil, errors.NewAccountNotFound(accountId)
	}
	if find.Error != nil {
		return nil, find.Error
	}
	return &existingUser, nil
}

func (r *AccountRepositoryImpl) Save(newAccount *model.Account) error {
	tx := r.Connection.Begin()
	tx.Create(newAccount)
	tx.Commit()
	return nil
}

func (r *AccountRepositoryImpl) Update(updatedAccount *model.Account) (*model.Account, error) {
	tx := r.Connection.Begin()
	tx.Save(updatedAccount)
	tx.Commit()
	return updatedAccount, nil
}

func NewAccountRepository(dbConnection *gorm.DB) AccountRepository {
	return &AccountRepositoryImpl{
		Connection: dbConnection,
	}
}
