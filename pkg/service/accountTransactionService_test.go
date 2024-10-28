package service

import (
	"context"
	"testing"
	"time"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/mrth1995/go-mockva/pkg/model"
	mockRepository "github.com/mrth1995/go-mockva/pkg/repository/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAccountTransactionService_Transfer(t *testing.T) {
	ctx := context.Background()

	initialSrcBalance := float64(1000_000)
	accountSrc := getAccountBalance(getAccountSrc(), initialSrcBalance)

	initialDstBalance := float64(200_000)
	accountDst := getAccountBalance(getAccountDst(), initialDstBalance)

	accountRepository := &mockRepository.MockAccountRepository{}
	accountTrxRepository := &mockRepository.MockAccountTransactionRepository{}

	accountRepository.On("FindAndLockAccountBalance", mock.Anything, accountSrc.ID).Return(accountSrc, nil)
	accountRepository.On("FindAndLockAccountBalance", mock.Anything, accountDst.ID).Return(accountDst, nil)
	var updatedAccountSrc *domain.AccountBalance
	var updatedAccountDst *domain.AccountBalance
	accountRepository.On("UpdateBalance", mock.Anything, accountSrc).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountSrc = args.Get(0).(*domain.AccountBalance)
	})
	accountRepository.On("UpdateBalance", mock.Anything, accountDst).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountDst = args.Get(0).(*domain.AccountBalance)
	})
	accountTrxRepository.On("Save", mock.Anything).Return(nil)

	accountService := NewAccountService(accountRepository)
	accountTrxService := NewAccountTrxService(accountService, accountTrxRepository)

	accountFundTransfer := &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       float64(100_000),
	}

	accountTransaction, err := accountTrxService.Transfer(ctx, accountFundTransfer)
	assertions := require.New(t)
	assertions.Nil(err, "No error")
	assertions.NotNil(accountTransaction, "Account transaction created")
	assertions.Equalf(accountSrc.ID, accountTransaction.AccountSrc.ID, "Source account should same")
	assertions.Equalf(accountDst.ID, accountTransaction.AccountDst.ID, "Destination account should same")
	assertions.Equalf(accountFundTransfer.Amount, accountTransaction.Amount, "Amount should same")
	assertions.Equalf(initialSrcBalance-accountTransaction.Amount, updatedAccountSrc.Balance, "Account src balance should same")
	assertions.Equalf(initialDstBalance+accountTransaction.Amount, updatedAccountDst.Balance, "Account dst balance should same")
}

func TestAccountTransactionService_TransferNegativeBalance(t *testing.T) {
	// accountSrc := getAccountSrcWithAllowNegativeBalance()
	// initialSrcBalance := accountSrc.Balance
	// accountDst := getAccountDst()
	// initialDstBalance := accountDst.Balance

	// accountRepository := &accountMock.MockAccountRepository{}
	// accountRepository.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	// accountRepository.On("FindById", accountDst.ID).Return(accountDst, nil)
	// var updatedAccountSrc *accountModel.Account
	// var updatedAccountDst *accountModel.Account
	// accountRepository.On("Update", accountSrc).Return(nil, nil).Run(func(args mock.Arguments) {
	// 	updatedAccountSrc = args.Get(0).(*accountModel.Account)
	// })
	// accountRepository.On("Update", accountDst).Return(nil, nil).Run(func(args mock.Arguments) {
	// 	updatedAccountDst = args.Get(0).(*accountModel.Account)
	// })
	// accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	// accountTrxRepository.On("Save", mock.Anything).Return(nil)
	// accountTrxService := &AccountTransactionServiceImpl{
	// 	AccountRepository:    accountRepository,
	// 	AccountTrxRepository: accountTrxRepository,
	// }
	// accountFundTransfer := &model.AccountFundTransfer{
	// 	AccountDstId: accountDst.ID,
	// 	AccountSrcId: accountSrc.ID,
	// 	Amount:       5000,
	// }
	// accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	// assertions := require.New(t)
	// assertions.Nil(err)
	// assertions.NotNil(accountTransaction, "Account transaction created")
	// assertions.Equalf(accountSrc.ID, accountTransaction.AccountSrc.ID, "Source account should same")
	// assertions.Equalf(accountDst.ID, accountTransaction.AccountDst.ID, "Destination account should same")
	// assertions.Equalf(accountFundTransfer.Amount, accountTransaction.Amount, "Amount should same")
	// assertions.Equalf(initialSrcBalance-accountTransaction.Amount, updatedAccountSrc.Balance, "Account src balance should same")
	// assertions.Equalf(initialDstBalance+accountTransaction.Amount, updatedAccountDst.Balance, "Account dst balance should same")
}

func TestAccountTransactionService_TransferWithSameAccount(t *testing.T) {
	// accountSrc := getAccountSrc()
	// accountDst := getAccountSrc()
	// accountRepository := &accountMock.MockAccountRepository{}
	// accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	// accountTrxService := &AccountTransactionServiceImpl{
	// 	AccountRepository:    accountRepository,
	// 	AccountTrxRepository: accountTrxRepository,
	// }
	// accountFundTransfer := &model.AccountFundTransfer{
	// 	AccountDstId: accountDst.ID,
	// 	AccountSrcId: accountSrc.ID,
	// 	Amount:       5000,
	// }
	// accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	// assertions := require.New(t)
	// assertions.Nil(accountTransaction)
	// assertions.NotNil(err)
}

func TestAccountTransactionService_TransferWithNilAccountSrc(t *testing.T) {
	// accountDst := getAccountDst()
	// accountRepository := &accountMock.MockAccountRepository{}
	// accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	// accountTrxService := &AccountTransactionServiceImpl{
	// 	AccountRepository:    accountRepository,
	// 	AccountTrxRepository: accountTrxRepository,
	// }
	// accountFundTransfer := &model.AccountFundTransfer{
	// 	AccountDstId: accountDst.ID,
	// 	AccountSrcId: "",
	// 	Amount:       5000,
	// }
	// accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	// assertions := require.New(t)
	// assertions.Nil(accountTransaction)
	// assertions.NotNil(err)
}

func TestAccountTransactionService_TransferWithNilAccountDst(t *testing.T) {
	// accountSrc := getAccountSrc()
	// accountRepository := &accountMock.MockAccountRepository{}
	// accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	// accountTrxService := &AccountTransactionServiceImpl{
	// 	AccountRepository:    accountRepository,
	// 	AccountTrxRepository: accountTrxRepository,
	// }
	// accountFundTransfer := &model.AccountFundTransfer{
	// 	AccountDstId: "",
	// 	AccountSrcId: accountSrc.ID,
	// 	Amount:       5000,
	// }
	// accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	// assertions := require.New(t)
	// assertions.Nil(accountTransaction)
	// assertions.NotNil(err)
}

func TestAccountTransactionService_TransferWithNegativeBalance(t *testing.T) {
	// accountSrc := getAccountSrc()
	// accountDst := getAccountSrc()
	// accountRepository := &accountMock.MockAccountRepository{}
	// accountTrxRepository := &accountTrxMock.MockAccountTransactionRepository{}
	// accountTrxService := &AccountTransactionServiceImpl{
	// 	AccountRepository:    accountRepository,
	// 	AccountTrxRepository: accountTrxRepository,
	// }
	// accountFundTransfer := &model.AccountFundTransfer{
	// 	AccountDstId: accountDst.ID,
	// 	AccountSrcId: accountSrc.ID,
	// 	Amount:       -5000,
	// }
	// accountTransaction, err := accountTrxService.Transfer(accountFundTransfer)
	// assertions := require.New(t)
	// assertions.Nil(accountTransaction)
	// assertions.NotNil(err, "Invalid balance")
}

func TestAccountTransactionServiceImpl_Transfer_InsufficientFunds(t *testing.T) {
	// accountSrc := getAccountSrc()
	// accountDst := getAccountDst()
	// balance := float64(1000_000)
	// accountRepo := &accountMock.MockAccountRepository{}
	// accountRepo.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	// accountRepo.On("FindById", accountDst.ID).Return(accountDst, nil)
	// accountTrxRepo := &accountTrxMock.MockAccountTransactionRepository{}
	// service := &AccountTransactionServiceImpl{accountRepo, accountTrxRepo}
	// transaction, err := service.Transfer(&model.AccountFundTransfer{
	// 	AccountDstId: accountDst.ID,
	// 	AccountSrcId: accountSrc.ID,
	// 	Amount:       balance,
	// })
	// assertions := require.New(t)
	// assertions.Nil(transaction)
	// assertions.NotNilf(err, "Insufficient funds")
}

func TestAccountTransactionServiceImpl_Transfer_AccountSrcNotFound(t *testing.T) {
	// accountSrc := getAccountSrc()
	// accountDst := getAccountDst()
	// balance := float64(1000_000)
	// accountRepo := &accountMock.MockAccountRepository{}
	// accountRepo.On("FindById", accountSrc.ID).Return(nil, errors.NewAccountNotFound(accountSrc.ID))
	// accountRepo.On("FindById", accountDst.ID).Return(accountDst)
	// accountTrxRepo := &accountTrxMock.MockAccountTransactionRepository{}
	// service := &AccountTransactionServiceImpl{accountRepo, accountTrxRepo}
	// transaction, err := service.Transfer(&model.AccountFundTransfer{
	// 	AccountDstId: accountDst.ID,
	// 	AccountSrcId: accountSrc.ID,
	// 	Amount:       balance,
	// })
	// assertions := require.New(t)
	// assertions.Nil(transaction)
	// assertions.NotNilf(err, "account src not found")
}

func TestAccountTransactionServiceImpl_Transfer_AccountDstNotFound(t *testing.T) {
	// accountSrc := getAccountSrc()
	// accountDst := getAccountDst()
	// balance := float64(1000_000)
	// accountRepo := &accountMock.MockAccountRepository{}
	// accountRepo.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	// accountRepo.On("FindById", accountDst.ID).Return(nil, errors.NewAccountNotFound(accountDst.ID))
	// accountTrxRepo := &accountTrxMock.MockAccountTransactionRepository{}
	// service := &AccountTransactionServiceImpl{accountRepo, accountTrxRepo}
	// transaction, err := service.Transfer(&model.AccountFundTransfer{
	// 	AccountDstId: accountDst.ID,
	// 	AccountSrcId: accountSrc.ID,
	// 	Amount:       balance,
	// })
	// assertions := require.New(t)
	// assertions.Nil(transaction)
	// assertions.NotNilf(err, "Account dst not found")
}

func getAccountSrc() *domain.Account {
	addr := "Jl sadarmanah"
	birthDate, err := time.Parse(time.DateOnly, "1995-03-01")
	if err != nil {
		panic(err)
	}
	return &domain.Account{
		ID:        "001",
		AccountID: "001",
		Name:      "Account 001",
		Address:   addr,
		BirthDate: birthDate,
		Gender:    true,
	}
}

func getAccountDst() *domain.Account {
	addr := "Jl sadarmanah"
	birthDate, err := time.Parse(time.DateOnly, "1995-03-01")
	if err != nil {
		panic(err)
	}
	return &domain.Account{
		ID:        "002",
		AccountID: "002",
		Name:      "Account 002",
		Address:   addr,
		BirthDate: birthDate,
		Gender:    true,
	}
}

func getAccountBalance(account *domain.Account, balance float64) *domain.AccountBalance {
	return &domain.AccountBalance{
		ID:        account.ID,
		AccountID: account.AccountID,
		Balance:   balance,
		Account:   account,
	}
}
