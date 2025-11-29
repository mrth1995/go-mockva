package postgresql

import (
	"context"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/mrth1995/go-mockva/pkg/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepositoryImpl struct {
	Connection *gorm.DB
}

func NewAccountRepository(dbConnection *gorm.DB) repository.AccountRepository {
	return &AccountRepositoryImpl{
		Connection: dbConnection,
	}
}

func (r *AccountRepositoryImpl) FindByID(ctx context.Context, accountId string) (*domain.Account, error) {
	var existingUser domain.Account
	find := r.Connection.First(&existingUser, "id = ?", accountId)
	if find.Error != nil && find.Error == gorm.ErrRecordNotFound {
		return nil, errors.NewAccountNotFound(accountId)
	}
	if find.Error != nil {
		return nil, find.Error
	}
	return &existingUser, nil
}

func (r *AccountRepositoryImpl) Save(ctx context.Context, newAccount *domain.Account) error {
	tx := r.Connection.Begin()
	tx.Create(newAccount)
	tx.Commit()
	return nil
}

func (r *AccountRepositoryImpl) Update(ctx context.Context, updatedAccount *domain.Account) (*domain.Account, error) {
	tx := r.Connection.Begin()
	tx.Save(updatedAccount)
	tx.Commit()
	return updatedAccount, nil
}

func (r *AccountRepositoryImpl) FindAndLockAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error) {
	var existingAccountBalance domain.AccountBalance
	find := r.Connection.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existingAccountBalance, "id = ?", accountID)
	if find.Error != nil && find.Error == gorm.ErrRecordNotFound {
		return nil, errors.NewAccountNotFound(accountID)
	}
	if find.Error != nil {
		return nil, find.Error
	}
	return &existingAccountBalance, nil
}

// UpdateBalance updates the account balance within the provided transaction context.
// Parameters:
//   - ctx: The request context
//   - accountBalance: The AccountBalance to update
//   - tx: The GORM transaction context
//
// Returns:
//   - *domain.AccountBalance: The updated balance
//   - error: If the operation fails
func (r *AccountRepositoryImpl) UpdateBalance(ctx context.Context, accountBalance *domain.AccountBalance, tx *gorm.DB) (*domain.AccountBalance, error) {
	if err := tx.Save(accountBalance).Error; err != nil {
		return nil, err
	}
	return accountBalance, nil
}
