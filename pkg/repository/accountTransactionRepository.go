package repository

//go:generate mockgen -destination=mock/mockAccountTransactionRepository.go -package=mock github.com/mrth1995/go-mockva/pkg/repository AccountTransactionRepository

import (
	"github.com/mrth1995/go-mockva/pkg/domain"
	"gorm.io/gorm"
)

type AccountTransactionRepository interface {
	// Save persists an AccountTransaction within the provided transaction context.
	// Parameters:
	//   - trx: The AccountTransaction to save
	//   - tx: The GORM transaction context
	// Returns:
	//   - error: If the operation fails
	Save(trx *domain.AccountTransaction, tx *gorm.DB) error
}
