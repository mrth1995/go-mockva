package service

import (
	"testing"
	"time"

	"github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	accountMock "github.com/mrth1995/go-mockva/pkg/account/mock"
	accountModel "github.com/mrth1995/go-mockva/pkg/account/model"
	accountTrxMock "github.com/mrth1995/go-mockva/pkg/accounttransaction/mock"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/model"
)

func TestAccountTransactionService_Transfer(t *testing.T) {
	accountSrc := getAccountSrc()
	initialSrcBalance := accountSrc.Balance
	accountDst := getAccountDst()
	initialDstBalance := accountDst.Balance
	accountFundTransfer := &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       5000,
	}
	accountRepository := &accountMock.MockAccountRepository{}
	accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}

	accountRepository.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	accountRepository.On("FindById", accountDst.ID).Return(accountDst, nil)
	var updatedAccountSrc *accountModel.Account
	var updatedAccountDst *accountModel.Account
	accountRepository.On("Update", accountSrc).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountSrc = args.Get(0).(*accountModel.Account)
	})
	accountRepository.On("Update", accountDst).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountDst = args.Get(0).(*accountModel.Account)
	})
	accountTrxRepository.On("Save", mock.Anything).Return(nil)
	accountTrxService := &AccountTransactionServiceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}

	accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	assertions := require.New(t)
	assertions.Nil(err, "No error")
	assertions.NotNil(accountTransaction, "Account transaction created")
	assertions.Equalf(accountSrc.ID, accountTransaction.AccountSrc.ID, "Source account should same")
	assertions.Equalf(accountDst.ID, accountTransaction.AccountDst.ID, "Destination account should same")
	assertions.Equalf(accountFundTransfer.Amount, accountTransaction.Amount, "Amount should same")
	assertions.Equalf(initialSrcBalance-accountTransaction.Amount, updatedAccountSrc.Balance, "Account src amount should same")
	assertions.Equalf(initialDstBalance+accountTransaction.Amount, updatedAccountDst.Balance, "Account dst amount should same")
}

func TestAccountTransactionService_TransferNegativeBalance(t *testing.T) {
	accountSrc := getAccountSrcWithAllowNegativeBalance()
	initialSrcBalance := accountSrc.Balance
	accountDst := getAccountDst()
	initialDstBalance := accountDst.Balance

	accountRepository := &accountMock.MockAccountRepository{}
	accountRepository.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	accountRepository.On("FindById", accountDst.ID).Return(accountDst, nil)
	var updatedAccountSrc *accountModel.Account
	var updatedAccountDst *accountModel.Account
	accountRepository.On("Update", accountSrc).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountSrc = args.Get(0).(*accountModel.Account)
	})
	accountRepository.On("Update", accountDst).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountDst = args.Get(0).(*accountModel.Account)
	})
	accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	accountTrxRepository.On("Save", mock.Anything).Return(nil)
	accountTrxService := &AccountTransactionServiceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       5000,
	}
	accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	assertions := require.New(t)
	assertions.Nil(err)
	assertions.NotNil(accountTransaction, "Account transaction created")
	assertions.Equalf(accountSrc.ID, accountTransaction.AccountSrc.ID, "Source account should same")
	assertions.Equalf(accountDst.ID, accountTransaction.AccountDst.ID, "Destination account should same")
	assertions.Equalf(accountFundTransfer.Amount, accountTransaction.Amount, "Amount should same")
	assertions.Equalf(initialSrcBalance-accountTransaction.Amount, updatedAccountSrc.Balance, "Account src amount should same")
	assertions.Equalf(initialDstBalance+accountTransaction.Amount, updatedAccountDst.Balance, "Account dst amount should same")
}

func TestAccountTransactionService_TransferWithSameAccount(t *testing.T) {
	accountSrc := getAccountSrc()
	accountDst := getAccountSrc()
	accountRepository := &accountMock.MockAccountRepository{}
	accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	accountTrxService := &AccountTransactionServiceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       5000,
	}
	accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	assertions := require.New(t)
	assertions.Nil(accountTransaction)
	assertions.NotNil(err)
}

func TestAccountTransactionService_TransferWithNilAccountSrc(t *testing.T) {
	accountDst := getAccountDst()
	accountRepository := &accountMock.MockAccountRepository{}
	accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	accountTrxService := &AccountTransactionServiceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: "",
		Amount:       5000,
	}
	accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	assertions := require.New(t)
	assertions.Nil(accountTransaction)
	assertions.NotNil(err)
}

func TestAccountTransactionService_TransferWithNilAccountDst(t *testing.T) {
	accountSrc := getAccountSrc()
	accountRepository := &accountMock.MockAccountRepository{}
	accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	accountTrxService := &AccountTransactionServiceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &model.AccountFundTransfer{
		AccountDstId: "",
		AccountSrcId: accountSrc.ID,
		Amount:       5000,
	}
	accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	assertions := require.New(t)
	assertions.Nil(accountTransaction)
	assertions.NotNil(err)
}

func TestAccountTransactionService_TransferWithNegativeBalance(t *testing.T) {
	accountSrc := getAccountSrc()
	accountDst := getAccountSrc()
	accountRepository := &accountMock.MockAccountRepository{}
	accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	accountTrxService := &AccountTransactionServiceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       -5000,
	}
	accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	assertions := require.New(t)
	assertions.Nil(accountTransaction)
	assertions.NotNil(err, "Invalid amount")
}

func TestAccountTransactionServiceImpl_Transfer_InsufficientFunds(t *testing.T) {
	accountSrc := getAccountSrc()
	accountDst := getAccountDst()
	amount := float64(1000_000)
	accountRepo := &accountMock.MockAccountRepository{}
	accountRepo.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	accountRepo.On("FindById", accountDst.ID).Return(accountDst, nil)
	accountTrxRepo := &accountTrxMock.MockAccountTransactionRepository{}
	service := &AccountTransactionServiceImpl{accountRepo, accountTrxRepo}
	transaction, err := service.Transfer(&model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       amount,
	})
	assertions := require.New(t)
	assertions.Nil(transaction)
	assertions.NotNilf(err, "Insufficient funds")
}

func TestAccountTransactionServiceImpl_Transfer_AccountSrcNotFound(t *testing.T) {
	accountSrc := getAccountSrc()
	accountDst := getAccountDst()
	amount := float64(1000_000)
	accountRepo := &accountMock.MockAccountRepository{}
	accountRepo.On("FindById", accountSrc.ID).Return(nil, errors.NewAccountNotFound(accountSrc.ID))
	accountRepo.On("FindById", accountDst.ID).Return(accountDst)
	accountTrxRepo := &accountTrxMock.MockAccountTransactionRepository{}
	service := &AccountTransactionServiceImpl{accountRepo, accountTrxRepo}
	transaction, err := service.Transfer(&model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       amount,
	})
	assertions := require.New(t)
	assertions.Nil(transaction)
	assertions.NotNilf(err, "account src not found")
}

func TestAccountTransactionServiceImpl_Transfer_AccountDstNotFound(t *testing.T) {
	accountSrc := getAccountSrc()
	accountDst := getAccountDst()
	amount := float64(1000_000)
	accountRepo := &accountMock.MockAccountRepository{}
	accountRepo.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	accountRepo.On("FindById", accountDst.ID).Return(nil, errors.NewAccountNotFound(accountDst.ID))
	accountTrxRepo := &accountTrxMock.MockAccountTransactionRepository{}
	service := &AccountTransactionServiceImpl{accountRepo, accountTrxRepo}
	transaction, err := service.Transfer(&model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       amount,
	})
	assertions := require.New(t)
	assertions.Nil(transaction)
	assertions.NotNilf(err, "Account dst not found")
}

func getAccountSrc() *accountModel.Account {
	addr := "Jl sadarmanah"
	return &accountModel.Account{
		ID:                   "001",
		Name:                 "Account 001",
		Address:              addr,
		BirthDate:            &time.Time{},
		Gender:               true,
		AllowNegativeBalance: false,
		Balance:              100_000,
	}
}

func getAccountDst() *accountModel.Account {
	addr := "Jl sadarmanah"
	return &accountModel.Account{
		ID:                   "002",
		Name:                 "Account 002",
		Address:              addr,
		BirthDate:            &time.Time{},
		Gender:               true,
		AllowNegativeBalance: false,
		Balance:              0,
	}
}

func getAccountSrcWithAllowNegativeBalance() *accountModel.Account {
	addr := "Jl sadarmanah"
	return &accountModel.Account{
		ID:                   "001",
		Name:                 "Account 001",
		Address:              addr,
		BirthDate:            &time.Time{},
		Gender:               true,
		AllowNegativeBalance: true,
		Balance:              0,
	}
}
