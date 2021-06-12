package repository

import (
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/model"
	"gorm.io/gorm"
)

type AccountTrxRepositoryImpl struct {
	Connection *gorm.DB
}

func (r *AccountTrxRepositoryImpl) Save(trx *model.AccountTransaction) error {
	tx := r.Connection.Begin()
	tx.Create(trx)
	tx.Commit()
	return nil
}

func NewAccountTrxRepository(dbConnection *gorm.DB) AccountTransactionRepository {
	return &AccountTrxRepositoryImpl{
		Connection: dbConnection,
	}
}
