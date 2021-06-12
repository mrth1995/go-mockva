package service

import (
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/model"
)

type AccountTransactionService interface {
	Transfer(accountFundTransfer *model.AccountFundTransfer) (*model.AccountTransaction, error)
}
