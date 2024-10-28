package postgresql

import (
	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/mrth1995/go-mockva/pkg/repository"
	"gorm.io/gorm"
)

type AccountTrxRepositoryImpl struct {
	Connection *gorm.DB
}

func NewAccountTrxRepository(dbConnection *gorm.DB) repository.AccountTransactionRepository {
	return &AccountTrxRepositoryImpl{
		Connection: dbConnection,
	}
}

func (r *AccountTrxRepositoryImpl) Save(trx *domain.AccountTransaction) error {
	tx := r.Connection.Begin()
	tx.Create(trx)
	tx.Commit()
	return nil
}
