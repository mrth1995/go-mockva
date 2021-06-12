package service

import (
	"errors"
	"strconv"
	"time"

	accountRepository "github.com/mrth1995/go-mockva/pkg/account/repository"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/model"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/repository"
)

type AccountTransactionServiceImpl struct {
	AccountRepository    accountRepository.AccountRepository
	AccountTrxRepository repository.AccountTransactionRepository
}

func (s *AccountTransactionServiceImpl) Transfer(accountFundTransfer *model.AccountFundTransfer) (*model.AccountTransaction, error) {
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

	accountSrc, err := s.AccountRepository.FindById(accountFundTransfer.AccountSrcId)
	if err != nil {
		return nil, err
	}
	accountDst, err := s.AccountRepository.FindById(accountFundTransfer.AccountDstId)
	if err != nil {
		return nil, err
	}
	if accountSrc.Balance-accountFundTransfer.Amount < 0 && !accountSrc.AllowNegativeBalance {
		return nil, errors.New("insufficient amount")
	}
	trxId := accountSrc.ID + ":" + accountDst.ID + ":" + strconv.Itoa(time.Now().Nanosecond())
	accountTrx := &model.AccountTransaction{
		ID:                   trxId,
		TransactionTimestamp: time.Now(),
		Amount:               accountFundTransfer.Amount,
		AccountSrc:           accountSrc,
		AccountDst:           accountDst,
	}
	accountSrc.Balance = accountSrc.Balance - accountFundTransfer.Amount
	accountDst.Balance = accountDst.Balance + accountFundTransfer.Amount

	err = s.AccountTrxRepository.Save(accountTrx)
	if err != nil {
		return nil, err
	}
	_, err = s.AccountRepository.Update(accountSrc)
	if err != nil {
		return nil, err
	}
	_, err = s.AccountRepository.Update(accountDst)
	if err != nil {
		return nil, err
	}

	return accountTrx, nil
}
