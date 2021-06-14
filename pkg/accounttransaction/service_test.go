package accounttransaction

import (
	"testing"
	"time"

	"github.com/mrth1995/go-mockva/pkg/account"

	"github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
)

func TestAccountTransactionService_Transfer(t *testing.T) {
	accountSrc := getAccountSrc()
	initialSrcBalance := accountSrc.Balance
	accountDst := getAccountDst()
	initialDstBalance := accountDst.Balance
	accountFundTransfer := &AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       5000,
	}
	accountRepository := &account.MockRepository{}
	accountTrxRepository := &MockRepository{}

	accountRepository.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	accountRepository.On("FindById", accountDst.ID).Return(accountDst, nil)
	var updatedAccountSrc *account.Account
	var updatedAccountDst *account.Account
	accountRepository.On("Update", accountSrc).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountSrc = args.Get(0).(*account.Account)
	})
	accountRepository.On("Update", accountDst).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountDst = args.Get(0).(*account.Account)
	})
	accountTrxRepository.On("Save", mock.Anything).Return(nil)
	accountTrxService := &serviceImpl{
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

	accountRepository := &account.MockRepository{}
	accountRepository.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	accountRepository.On("FindById", accountDst.ID).Return(accountDst, nil)
	var updatedAccountSrc *account.Account
	var updatedAccountDst *account.Account
	accountRepository.On("Update", accountSrc).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountSrc = args.Get(0).(*account.Account)
	})
	accountRepository.On("Update", accountDst).Return(nil, nil).Run(func(args mock.Arguments) {
		updatedAccountDst = args.Get(0).(*account.Account)
	})
	accountTrxRepository := &MockRepository{}
	accountTrxRepository.On("Save", mock.Anything).Return(nil)
	accountTrxService := &serviceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &AccountFundTransfer{
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
	accountRepository := &account.MockRepository{}
	accountTrxRepository := &MockRepository{}
	accountTrxService := &serviceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &AccountFundTransfer{
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
	accountRepository := &account.MockRepository{}
	accountTrxRepository := &MockRepository{}
	accountTrxService := &serviceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &AccountFundTransfer{
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
	accountRepository := &account.MockRepository{}
	accountTrxRepository := &MockRepository{}
	accountTrxService := &serviceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &AccountFundTransfer{
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
	accountRepository := &account.MockRepository{}
	accountTrxRepository := &MockRepository{}
	accountTrxService := &serviceImpl{
		AccountRepository:    accountRepository,
		AccountTrxRepository: accountTrxRepository,
	}
	accountFundTransfer := &AccountFundTransfer{
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
	accountRepo := &account.MockRepository{}
	accountRepo.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	accountRepo.On("FindById", accountDst.ID).Return(accountDst, nil)
	accountTrxRepo := &MockRepository{}
	service := &serviceImpl{accountRepo, accountTrxRepo}
	transaction, err := service.Transfer(&AccountFundTransfer{
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
	accountRepo := &account.MockRepository{}
	accountRepo.On("FindById", accountSrc.ID).Return(nil, errors.NewAccountNotFound(accountSrc.ID))
	accountRepo.On("FindById", accountDst.ID).Return(accountDst)
	accountTrxRepo := &MockRepository{}
	service := &serviceImpl{accountRepo, accountTrxRepo}
	transaction, err := service.Transfer(&AccountFundTransfer{
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
	accountRepo := &account.MockRepository{}
	accountRepo.On("FindById", accountSrc.ID).Return(accountSrc, nil)
	accountRepo.On("FindById", accountDst.ID).Return(nil, errors.NewAccountNotFound(accountDst.ID))
	accountTrxRepo := &MockRepository{}
	service := &serviceImpl{accountRepo, accountTrxRepo}
	transaction, err := service.Transfer(&AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       amount,
	})
	assertions := require.New(t)
	assertions.Nil(transaction)
	assertions.NotNilf(err, "Account dst not found")
}

func getAccountSrc() *account.Account {
	addr := "Jl sadarmanah"
	return &account.Account{
		ID:                   "001",
		Name:                 "Account 001",
		Address:              addr,
		BirthDate:            &time.Time{},
		Gender:               true,
		AllowNegativeBalance: false,
		Balance:              100_000,
	}
}

func getAccountDst() *account.Account {
	addr := "Jl sadarmanah"
	return &account.Account{
		ID:                   "002",
		Name:                 "Account 002",
		Address:              addr,
		BirthDate:            &time.Time{},
		Gender:               true,
		AllowNegativeBalance: false,
		Balance:              0,
	}
}

func getAccountSrcWithAllowNegativeBalance() *account.Account {
	addr := "Jl sadarmanah"
	return &account.Account{
		ID:                   "001",
		Name:                 "Account 001",
		Address:              addr,
		BirthDate:            &time.Time{},
		Gender:               true,
		AllowNegativeBalance: true,
		Balance:              0,
	}
}
