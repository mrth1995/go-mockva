package service

import (
	"context"
	"testing"
	"time"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/mrth1995/go-mockva/pkg/model"
	accountMock "github.com/mrth1995/go-mockva/pkg/repository/mock"
	"github.com/mrth1995/go-mockva/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	accountID = "12345"
)

func TestAccountServiceImpl_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	repository := accountMock.NewMockAccountRepository(ctrl)

	// Mock FindByID to return not found error
	repository.EXPECT().
		FindByID(gomock.Any(), "100").
		Return(nil, errors.NewAccountNotFound("100"))

	// Mock Save to capture the saved account
	var capturedAccount *domain.Account
	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, acc *domain.Account) error {
			capturedAccount = acc
			return nil
		})

	accountService := &AccountServiceImpl{accountRepository: repository}
	account, accountAlreadyExist := accountService.Register(ctx, accountRegister)

	assertions := require.New(t)
	assertions.Nil(accountAlreadyExist, "Should not error")
	assertions.NotNilf(account, "Created account should not be empty")
	assertions.Equalf(capturedAccount.ID, account.ID, "Account ID should equals")
	assertions.Equalf(capturedAccount.Address, account.Address, "Address should equals")
	assertions.Equalf(capturedAccount.BirthDate, account.BirthDate, "BirthDate should equals")
	assertions.Equalf(capturedAccount.Gender, account.Gender, "Gender should equals")
	assertions.Equalf(capturedAccount.Name, account.Name, "Name should equals")
}

func TestAccountServiceImpl_Register_AccountAlreadyExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	repository := accountMock.NewMockAccountRepository(ctrl)
	repository.EXPECT().
		FindByID(gomock.Any(), accountRegister.ID).
		Return(&domain.Account{}, nil)

	accountService := &AccountServiceImpl{
		accountRepository: repository,
	}

	newAccount, alreadyExist := accountService.Register(ctx, accountRegister)
	assertions := require.New(t)
	assertions.Nil(newAccount)
	assertions.NotNilf(alreadyExist, "Error already exist should not be nil")
}

func TestAccountServiceImpl_Edit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	repository := accountMock.NewMockAccountRepository(ctrl)
	repository.EXPECT().
		FindByID(gomock.Any(), accountID).
		Return(existingAccount, nil)

	var capturedAccount *domain.Account
	repository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, acc *domain.Account) (*domain.Account, error) {
			capturedAccount = acc
			return nil, nil
		})

	service := &AccountServiceImpl{accountRepository: repository}
	account, err := service.Edit(ctx, accountID, edit)
	assertions := require.New(t)
	assertions.Nil(err)
	assertions.NotNilf(account, "Updated account not nil")
	assertions.Equal(*edit.Name, capturedAccount.Name)
	assertions.Equal(*edit.Address, capturedAccount.Address)
	assertions.Equal(*edit.Gender, capturedAccount.Gender)
	assertions.Equal(editedBirthdate.String(), capturedAccount.BirthDate.String())
}

func TestAccountServiceImpl_Edit_AccountNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	edit := getEditAccount()
	repository := accountMock.NewMockAccountRepository(ctrl)
	repository.EXPECT().
		FindByID(gomock.Any(), accountID).
		Return(nil, errors.NewAccountNotFound(accountID))

	service := &AccountServiceImpl{
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
