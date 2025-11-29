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

// Save persists an AccountTransaction within the provided transaction context.
// Parameters:
//   - trx: The AccountTransaction to save
//   - tx: The GORM transaction context
//
// Returns:
//   - error: If the operation fails
func (r *AccountTrxRepositoryImpl) Save(trx *domain.AccountTransaction, tx *gorm.DB) error {
	if err := tx.Create(trx).Error; err != nil {
		return err
	}
	return nil
}
