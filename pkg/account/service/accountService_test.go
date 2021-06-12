package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/mrth1995/go-mockva/pkg/errors"

	accountMock "github.com/mrth1995/go-mockva/pkg/account/mock"

	"github.com/mrth1995/go-mockva/pkg/account/model"
	"github.com/stretchr/testify/require"
)

func TestAccountServiceImpl_Register(t *testing.T) {
	birthDate := "1996-03-11"
	accountRegister := &model.AccountRegister{
		ID:                   "100",
		Name:                 "Siska",
		Address:              "Jl sadarmanah",
		BirthDate:            birthDate,
		Gender:               false,
		AllowNegativeBalance: false,
	}
	repository := new(accountMock.MockAccountRepository)
	repository.On("FindById", "100").Return(nil, fmt.Errorf("account %v not found", accountRegister.ID))
	//birtDateTime, _ := time.Parse("2006-01-02", birthDate)
	var newAccount *model.Account
	repository.On("Save", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		newAccount = args.Get(0).(*model.Account)
	})
	accountService := &AccountServiceImpl{repository: repository}
	account, accountAlreadyExist := accountService.Register(accountRegister)

	assertions := require.New(t)
	assertions.Nil(accountAlreadyExist, "Should not error")
	assertions.NotNilf(account, "Created account should not be empty")
	assertions.Equalf(newAccount, account, "Should equals")
}

func TestAccountServiceImpl_Register_AccountAlreadyExist(t *testing.T) {
	dateString := "1996-03-11"
	accountRegister := &model.AccountRegister{
		ID:                   "111",
		Name:                 "Ridwan",
		Address:              "Jl sadarmanah",
		BirthDate:            dateString,
		Gender:               true,
		AllowNegativeBalance: false,
	}

	repository := &accountMock.MockAccountRepository{}
	repository.On("FindById", accountRegister.ID).Return(&model.Account{}, nil)

	accountService := &AccountServiceImpl{
		repository: repository,
	}

	newAccount, alreadyExist := accountService.Register(accountRegister)
	assertions := require.New(t)
	assertions.Nil(newAccount)
	assertions.NotNilf(alreadyExist, "Error already exist should not be nil")
}

func TestAccountServiceImpl_Edit(t *testing.T) {
	edit := getEditAccount()
	editedBirthdate, _ := time.Parse("2006-01-02", edit.BirthDate)
	existingBirthDate, _ := time.Parse("2006-01-02", "1995-08-25")
	addr := "Jl menteng atas"
	existingAccount := &model.Account{
		ID:                   "001",
		Name:                 "Ridwan Taufik",
		Address:              addr,
		BirthDate:            &existingBirthDate,
		Gender:               false,
		Balance:              0,
		AllowNegativeBalance: true,
	}

	repository := &accountMock.MockAccountRepository{}
	repository.On("FindById", edit.ID).Return(existingAccount, nil)
	var updatedAccount *model.Account
	repository.On("Update", mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccount = args.Get(0).(*model.Account)
	})
	service := &AccountServiceImpl{repository: repository}
	account, err := service.Edit(edit)
	assertions := require.New(t)
	assertions.Nil(err)
	assertions.NotNilf(account, "Updated account not nil")
	assertions.Equal(edit.Name, updatedAccount.Name)
	assertions.Equal(edit.Address, updatedAccount.Address)
	assertions.Equal(edit.Gender, updatedAccount.Gender)
	assertions.Equal(edit.AllowNegativeBalance, updatedAccount.AllowNegativeBalance)
	assertions.Equal(editedBirthdate.String(), updatedAccount.BirthDate.String())
}

func TestAccountServiceImpl_Edit_AccountNotFound(t *testing.T) {
	edit := getEditAccount()
	repository := &accountMock.MockAccountRepository{}
	repository.On("FindById", edit.ID).Return(nil, errors.NewAccountNotFound(edit.ID))
	service := &AccountServiceImpl{
		repository: repository,
	}
	account, err := service.Edit(edit)
	assertions := require.New(t)
	assertions.Nil(account)
	assertions.NotNil(err, "Account not found")
}

func getEditAccount() *model.AccountEdit {
	return &model.AccountEdit{
		ID:                   "001",
		Name:                 "Ridwan",
		Address:              "JL sadarmanah",
		BirthDate:            "1995-03-01",
		Gender:               true,
		AllowNegativeBalance: false,
	}
}
