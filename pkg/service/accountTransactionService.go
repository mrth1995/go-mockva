package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/mrth1995/go-mockva/pkg/model"
	"github.com/mrth1995/go-mockva/pkg/repository"
)

type AccountTransactionService struct {
	accountService       *AccountService
	accountTrxRepository repository.AccountTransactionRepository
}

func NewAccountTrxService(accountService *AccountService, accountTrxRepo repository.AccountTransactionRepository) *AccountTransactionService {
	return &AccountTransactionService{
		accountService:       accountService,
		accountTrxRepository: accountTrxRepo,
	}
}

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

	accountSrc, err := s.accountService.FindAndLockAccountBalance(ctx, accountFundTransfer.AccountSrcId)
	if err != nil {
		return nil, err
	}
	accountDst, err := s.accountService.FindAndLockAccountBalance(ctx, accountFundTransfer.AccountDstId)
	if err != nil {
		return nil, err
	}
	if accountSrc.Balance-accountFundTransfer.Amount < 0 && !accountSrc.AllowNegativeBalance {
		return nil, errors.New("insufficient amount")
	}
	trxId := accountSrc.ID + ":" + accountDst.ID + ":" + strconv.Itoa(time.Now().Nanosecond())
	accountTrx := &domain.AccountTransaction{
		ID:                   trxId,
		TransactionTimestamp: time.Now(),
		Amount:               accountFundTransfer.Amount,
		AccountSrc:           accountSrc,
		AccountDst:           accountDst,
	}
	accountSrc.Balance = accountSrc.Balance - accountFundTransfer.Amount
	accountDst.Balance = accountDst.Balance + accountFundTransfer.Amount

	err = s.accountTrxRepository.Save(accountTrx)
	if err != nil {
		return nil, err
	}
	_, err = s.accountService.UpdateBalance(ctx, accountSrc)
	if err != nil {
		return nil, err
	}
	_, err = s.accountService.UpdateBalance(ctx, accountDst)
	if err != nil {
		return nil, err
	}
	return accountTrx, nil
}
