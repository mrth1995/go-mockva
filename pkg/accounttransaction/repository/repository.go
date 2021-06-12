package repository

import "github.com/mrth1995/go-mockva/pkg/accounttransaction/model"

type AccountTransactionRepository interface {
	Save(trx *model.AccountTransaction) error
}
