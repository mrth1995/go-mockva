package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/mrth1995/go-mockva/pkg/model"
	"github.com/mrth1995/go-mockva/pkg/repository"
	"gorm.io/gorm"
)

type AccountTransactionService struct {
	accountService       AccountService
	accountTrxRepository repository.AccountTransactionRepository
	txManager            repository.DBTransactionManager
}

func NewAccountTrxService(accountService AccountService, accountTrxRepo repository.AccountTransactionRepository, txManager repository.DBTransactionManager) *AccountTransactionService {
	return &AccountTransactionService{
		accountService:       accountService,
		accountTrxRepository: accountTrxRepo,
		txManager:            txManager,
	}
}

// Transfer moves funds between two accounts atomically using a database transaction.
// Parameters:
//   - ctx: The request context
//   - accountFundTransfer: The transfer details
//
// Returns:
//   - *domain.AccountTransaction: The completed transaction
//   - error: If the operation fails
func (s *AccountTransactionService) Transfer(ctx context.Context, accountFundTransfer *model.AccountFundTransfer) (*domain.AccountTransaction, error) {
	if accountFundTransfer.AccountSrcId == "" {
		return nil, errors.New("account src cannot be empty")
	}
	if accountFundTransfer.AccountDstId == "" {
		return nil, errors.New("account dst cannot be empty")
	}
	if accountFundTransfer.Amount <= 0 {
		return nil, errors.New("invalid amount")
	}
	if accountFundTransfer.AccountDstId == accountFundTransfer.AccountSrcId {
		return nil, errors.New("cannot transfer with same account")
	}

	var accountTrx *domain.AccountTransaction
	err := s.txManager.Transaction(func(tx *gorm.DB) error {
		accountSrc, err := s.accountService.FindAndLockAccountBalance(ctx, accountFundTransfer.AccountSrcId)
		if err != nil {
			return err
		}
		accountDst, err := s.accountService.FindAndLockAccountBalance(ctx, accountFundTransfer.AccountDstId)
		if err != nil {
			return err
		}
		if accountSrc.Balance-accountFundTransfer.Amount < 0 && !accountSrc.AllowNegativeBalance {
			return errors.New("insufficient amount")
		}
		trxId := accountSrc.ID + ":" + accountDst.ID + ":" + strconv.Itoa(time.Now().Nanosecond())
		accountTrx = &domain.AccountTransaction{
			ID:                   trxId,
			TransactionTimestamp: time.Now(),
			Amount:               accountFundTransfer.Amount,
			AccountSrc:           accountSrc,
			AccountDst:           accountDst,
		}
		accountSrc.Balance = accountSrc.Balance - accountFundTransfer.Amount
		accountDst.Balance = accountDst.Balance + accountFundTransfer.Amount

		if err := s.accountTrxRepository.Save(accountTrx, tx); err != nil {
			return err
		}
		_, err = s.accountService.UpdateBalance(ctx, accountSrc, tx)
		if err != nil {
			return err
		}
		_, err = s.accountService.UpdateBalance(ctx, accountDst, tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return accountTrx, nil
}
