package repository

import "github.com/mrth1995/go-mockva/pkg/domain"

type AccountTransactionRepository interface {
	Save(trx *domain.AccountTransaction) error
}
