package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/mrth1995/go-mockva/pkg/utils"

	accountMock "github.com/mrth1995/go-mockva/pkg/repository/mock"

	"github.com/mrth1995/go-mockva/pkg/model"
	"github.com/stretchr/testify/require"
)

const (
	accountID = "12345"
)

func TestAccountServiceImpl_Register(t *testing.T) {
	ctx := context.Background()

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
	repository.On("FindByID", mock.Anything, "100").Return(nil, fmt.Errorf("account %v not found", accountRegister.ID))
	//birtDateTime, _ := time.Parse("2006-01-02", birthDate)
	var newAccount *domain.Account
	repository.On("Save", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		newAccount = args.Get(1).(*domain.Account)
	})
	accountService := &AccountService{accountRepository: repository}
	account, accountAlreadyExist := accountService.Register(ctx, accountRegister)

	assertions := require.New(t)
	assertions.Nil(accountAlreadyExist, "Should not error")
	assertions.NotNilf(account, "Created account should not be empty")
	assertions.Equalf(newAccount.AccountID, account.AccountID, "Account ID should equals")
	assertions.Equalf(newAccount.Address, account.Address, "Address should equals")
	assertions.Equalf(newAccount.BirthDate, account.BirthDate, "BirthDate should equals")
	assertions.Equalf(newAccount.Gender, account.Gender, "Gender should equals")
	assertions.Equalf(newAccount.Name, account.Name, "Name should equals")
}

func TestAccountServiceImpl_Register_AccountAlreadyExist(t *testing.T) {
	ctx := context.Background()

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
	repository.On("FindByID", mock.Anything, accountRegister.ID).Return(&domain.Account{}, nil)

	accountService := &AccountService{
		accountRepository: repository,
	}

	newAccount, alreadyExist := accountService.Register(ctx, accountRegister)
	assertions := require.New(t)
	assertions.Nil(newAccount)
	assertions.NotNilf(alreadyExist, "Error already exist should not be nil")
}

func TestAccountServiceImpl_Edit(t *testing.T) {
	ctx := context.Background()

	edit := getEditAccount()
	editedBirthdate, _ := time.Parse(time.DateOnly, *edit.BirthDate)
	existingBirthDate, _ := time.Parse(time.DateOnly, "1995-08-25")
	addr := "Jl menteng atas"
	existingAccount := &domain.Account{
		ID:        accountID,
		Name:      "Ridwan Taufik",
		Address:   addr,
		BirthDate: existingBirthDate,
		Gender:    false,
	}

	repository := &accountMock.MockAccountRepository{}
	repository.On("FindByID", mock.Anything, accountID).Return(existingAccount, nil)
	var updatedAccount *domain.Account
	repository.On("Update", mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccount = args.Get(1).(*domain.Account)
	})
	service := &AccountService{accountRepository: repository}
	account, err := service.Edit(ctx, accountID, edit)
	assertions := require.New(t)
	assertions.Nil(err)
	assertions.NotNilf(account, "Updated account not nil")
	assertions.Equal(*edit.Name, updatedAccount.Name)
	assertions.Equal(*edit.Address, updatedAccount.Address)
	assertions.Equal(*edit.Gender, updatedAccount.Gender)
	assertions.Equal(editedBirthdate.String(), updatedAccount.BirthDate.String())
}

func TestAccountServiceImpl_Edit_AccountNotFound(t *testing.T) {
	ctx := context.Background()

	edit := getEditAccount()
	repository := &accountMock.MockAccountRepository{}
	repository.On("FindByID", mock.Anything, accountID).Return(nil, errors.NewAccountNotFound(accountID))
	service := &AccountService{
		accountRepository: repository,
	}
	account, err := service.Edit(ctx, accountID, edit)
	assertions := require.New(t)
	assertions.Nil(account)
	assertions.NotNil(err, "Account not found")
}

func getEditAccount() *model.AccountEdit {
	return &model.AccountEdit{
		Name:      utils.ToStringPointer("Ridwan"),
		Address:   utils.ToStringPointer("JL sadarmanah"),
		BirthDate: utils.ToStringPointer("1995-03-01"),
		Gender:    utils.ToBooleanPointer(true),
	}
}
