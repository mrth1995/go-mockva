package service

import (
	"context"
	"testing"
	"time"

	"github.com/mrth1995/go-mockva/pkg/domain"
	pkgErrors "github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/mrth1995/go-mockva/pkg/model"
	mockRepo "github.com/mrth1995/go-mockva/pkg/repository/mock"
	mockService "github.com/mrth1995/go-mockva/pkg/service/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// Note: Full integration test for TestAccountTransactionService_Transfer requires a real database
// due to the service's direct dependency on the DB connection for transaction management.
// The following tests focus on validation logic and error paths that can be tested without DB transactions.

func TestAccountTransactionService_Transfer_ValidationErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	accountService := mockService.NewMockAccountService(ctrl)
	accountTrxService := &AccountTransactionService{
		accountService: accountService,
	}

	testCases := []struct {
		name        string
		transfer    *model.AccountFundTransfer
		expectedErr string
	}{
		{
			name: "Empty source account",
			transfer: &model.AccountFundTransfer{
				AccountSrcId: "",
				AccountDstId: "002",
				Amount:       100000,
			},
			expectedErr: "account src cannot be empty",
		},
		{
			name: "Empty destination account",
			transfer: &model.AccountFundTransfer{
				AccountSrcId: "001",
				AccountDstId: "",
				Amount:       100000,
			},
			expectedErr: "account dst cannot be empty",
		},
		{
			name: "Invalid amount (zero)",
			transfer: &model.AccountFundTransfer{
				AccountSrcId: "001",
				AccountDstId: "002",
				Amount:       0,
			},
			expectedErr: "invalid amount",
		},
		{
			name: "Invalid amount (negative)",
			transfer: &model.AccountFundTransfer{
				AccountSrcId: "001",
				AccountDstId: "002",
				Amount:       -100,
			},
			expectedErr: "invalid amount",
		},
		{
			name: "Same source and destination",
			transfer: &model.AccountFundTransfer{
				AccountSrcId: "001",
				AccountDstId: "001",
				Amount:       100000,
			},
			expectedErr: "cannot transfer with same account",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := accountTrxService.Transfer(ctx, tc.transfer)
			assertions := require.New(t)
			assertions.Nil(result)
			assertions.NotNil(err)
			assertions.Contains(err.Error(), tc.expectedErr)
		})
	}
}

// TestAccountTransactionService_Transfer tests successful fund transfer
func TestAccountTransactionService_Transfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	initialSrcBalance := float64(1_000_000)
	accountSrc := getAccountBalance(getAccountSrc(), initialSrcBalance)

	initialDstBalance := float64(200_000)
	accountDst := getAccountBalance(getAccountDst(), initialDstBalance)

	accountService := mockService.NewMockAccountService(ctrl)
	accountTrxRepo := mockRepo.NewMockAccountTransactionRepository(ctrl)
	txManager := mockRepo.NewMockDBTransactionManager(ctrl)

	txManager.EXPECT().
		Transaction(gomock.Any()).
		DoAndReturn(func(fc func(tx *gorm.DB) error) error {
			return fc(nil)
		})

	accountService.EXPECT().
		FindAndLockAccountBalance(ctx, accountSrc.ID).
		Return(accountSrc, nil)

	accountService.EXPECT().
		FindAndLockAccountBalance(ctx, accountDst.ID).
		Return(accountDst, nil)

	accountTrxRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	accountService.EXPECT().
		UpdateBalance(ctx, gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, bal *domain.AccountBalance, tx *gorm.DB) (*domain.AccountBalance, error) {
			require.Equal(t, initialSrcBalance-100_000, bal.Balance)
			return bal, nil
		})

	accountService.EXPECT().
		UpdateBalance(ctx, gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, bal *domain.AccountBalance, tx *gorm.DB) (*domain.AccountBalance, error) {
			require.Equal(t, initialDstBalance+100_000, bal.Balance)
			return bal, nil
		})

	accountTrxService := NewAccountTrxService(accountService, accountTrxRepo, txManager)

	accountFundTransfer := &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       100_000,
	}

	accountTransaction, err := accountTrxService.Transfer(ctx, accountFundTransfer)

	assertions := require.New(t)
	assertions.Nil(err, "No error")
	assertions.NotNil(accountTransaction, "Account transaction created")
	assertions.Equal(accountSrc.ID, accountTransaction.AccountSrc.ID, "Source account should match")
	assertions.Equal(accountDst.ID, accountTransaction.AccountDst.ID, "Destination account should match")
	assertions.Equal(accountFundTransfer.Amount, accountTransaction.Amount, "Amount should match")
}

func TestAccountTransactionService_TransferNegativeBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	initialSrcBalance := float64(100_000)
	accountSrc := getAccountSrcWithAllowNegativeBalance()
	accountSrc.Balance = initialSrcBalance

	initialDstBalance := float64(200_000)
	accountDst := getAccountBalance(getAccountDst(), initialDstBalance)

	accountService := mockService.NewMockAccountService(ctrl)
	accountTrxRepo := mockRepo.NewMockAccountTransactionRepository(ctrl)
	txManager := mockRepo.NewMockDBTransactionManager(ctrl)

	txManager.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fc func(tx *gorm.DB) error) error {
		return fc(nil)
	})

	accountService.EXPECT().FindAndLockAccountBalance(ctx, accountSrc.ID).Return(accountSrc, nil)
	accountService.EXPECT().FindAndLockAccountBalance(ctx, accountDst.ID).Return(accountDst, nil)
	accountTrxRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
	accountService.EXPECT().UpdateBalance(ctx, gomock.Any(), gomock.Any()).Return(accountSrc, nil)
	accountService.EXPECT().UpdateBalance(ctx, gomock.Any(), gomock.Any()).Return(accountDst, nil)

	accountTrxService := NewAccountTrxService(accountService, accountTrxRepo, txManager)

	accountFundTransfer := &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       150_000, // More than balance, but negative allowed
	}

	accountTransaction, err := accountTrxService.Transfer(ctx, accountFundTransfer)
	assertions := require.New(t)
	assertions.Nil(err)
	assertions.NotNil(accountTransaction, "Account transaction created")
	assertions.Equal(accountSrc.ID, accountTransaction.AccountSrc.ID)
	assertions.Equal(accountDst.ID, accountTransaction.AccountDst.ID)
	assertions.Equal(accountFundTransfer.Amount, accountTransaction.Amount)
}

func TestAccountTransactionServiceImpl_Transfer_InsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	initialBalance := float64(100_000)
	accountSrc := getAccountBalance(getAccountSrc(), initialBalance)
	accountDst := getAccountBalance(getAccountDst(), 200_000)

	accountService := mockService.NewMockAccountService(ctrl)
	accountTrxRepo := mockRepo.NewMockAccountTransactionRepository(ctrl)
	txManager := mockRepo.NewMockDBTransactionManager(ctrl)

	txManager.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fc func(tx *gorm.DB) error) error {
		return fc(nil)
	})

	accountService.EXPECT().FindAndLockAccountBalance(ctx, accountSrc.ID).Return(accountSrc, nil)
	accountService.EXPECT().FindAndLockAccountBalance(ctx, accountDst.ID).Return(accountDst, nil)

	accountTrxService := NewAccountTrxService(accountService, accountTrxRepo, txManager)

	transaction, err := accountTrxService.Transfer(ctx, &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       1_000_000, // More than available balance
	})

	assertions := require.New(t)
	assertions.Nil(transaction)
	assertions.NotNil(err, "Insufficient funds")
	assertions.Contains(err.Error(), "insufficient amount")
}

func TestAccountTransactionServiceImpl_Transfer_AccountSrcNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	accountSrc := getAccountSrc()
	accountDst := getAccountDst()

	accountService := mockService.NewMockAccountService(ctrl)
	accountTrxRepo := mockRepo.NewMockAccountTransactionRepository(ctrl)
	txManager := mockRepo.NewMockDBTransactionManager(ctrl)

	txManager.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fc func(tx *gorm.DB) error) error {
		return fc(nil)
	})

	accountService.EXPECT().
		FindAndLockAccountBalance(ctx, accountSrc.ID).
		Return(nil, pkgErrors.NewAccountNotFound(accountSrc.ID))

	accountTrxService := NewAccountTrxService(accountService, accountTrxRepo, txManager)

	transaction, err := accountTrxService.Transfer(ctx, &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       100_000,
	})

	assertions := require.New(t)
	assertions.Nil(transaction)
	assertions.NotNil(err, "account src not found")
}

func TestAccountTransactionServiceImpl_Transfer_AccountDstNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	initialBalance := float64(1_000_000)
	accountSrc := getAccountBalance(getAccountSrc(), initialBalance)
	accountDst := getAccountDst()

	accountService := mockService.NewMockAccountService(ctrl)
	accountTrxRepo := mockRepo.NewMockAccountTransactionRepository(ctrl)
	txManager := mockRepo.NewMockDBTransactionManager(ctrl)

	txManager.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fc func(tx *gorm.DB) error) error {
		return fc(nil)
	})

	accountService.EXPECT().FindAndLockAccountBalance(ctx, accountSrc.ID).Return(accountSrc, nil)
	accountService.EXPECT().
		FindAndLockAccountBalance(ctx, accountDst.ID).
		Return(nil, pkgErrors.NewAccountNotFound(accountDst.ID))

	accountTrxService := NewAccountTrxService(accountService, accountTrxRepo, txManager)

	transaction, err := accountTrxService.Transfer(ctx, &model.AccountFundTransfer{
		AccountDstId: accountDst.ID,
		AccountSrcId: accountSrc.ID,
		Amount:       100_000,
	})

	assertions := require.New(t)
	assertions.Nil(transaction)
	assertions.NotNil(err, "Account dst not found")
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

func getAccountSrcWithAllowNegativeBalance() *domain.AccountBalance {
	account := getAccountSrc()
	return &domain.AccountBalance{
		ID:                   account.ID,
		AccountID:            account.AccountID,
		Balance:              100_000,
		AllowNegativeBalance: true,
		Account:              account,
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
		ID:                   account.ID,
		AccountID:            account.AccountID,
		Balance:              balance,
		AllowNegativeBalance: false,
		Account:              account,
	}
}
