package repository

import (
	"github.com/mrth1995/go-mockva/pkg/account/model"
)

type AccountRepository interface {
	FindById(accountId string) (*model.Account, error)
	Save(newAccount *model.Account) error
	Update(updatedAccount *model.Account) (*model.Account, error)
}
