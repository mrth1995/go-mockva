package service

import (
	"github.com/mrth1995/go-mockva/pkg/account/model"
	"github.com/mrth1995/go-mockva/pkg/account/repository"
)

type AccountService interface {
	Register(register *model.AccountRegister) (*model.Account, error)
	Edit(id string, edit *model.AccountEdit) (*model.Account, error)
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return &AccountServiceImpl{repository: accountRepository}
}
