package accounttransaction

import (
	"strconv"
	"time"

	"github.com/mrth1995/go-mockva/pkg/account"
	"github.com/mrth1995/go-mockva/pkg/errors"
)

type Service interface {
	Transfer(accountFundTransfer *AccountFundTransfer) (*AccountTransaction, error)
}

type serviceImpl struct {
	AccountRepository    account.Repository
	AccountTrxRepository Repository
}

func (s *serviceImpl) Transfer(accountFundTransfer *AccountFundTransfer) (*AccountTransaction, error) {
	if accountFundTransfer.AccountSrcId == "" {
		return nil, errors.NewAccountNotFound("account src cannot be empty")
	}
	if accountFundTransfer.AccountDstId == "" {
		return nil, errors.NewAccountNotFound("account dst cannot be empty")
	}
	if accountFundTransfer.Amount <= 0 {
		return nil, errors.NewInvalidAmount()
	}
	if accountFundTransfer.AccountDstId == accountFundTransfer.AccountSrcId {
		return nil, errors.NewInvalidAccount()
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
		return nil, errors.NewInsufficientAmount()
	}
	trxId := accountSrc.ID + ":" + accountDst.ID + ":" + strconv.Itoa(time.Now().Nanosecond())
	accountTrx := &AccountTransaction{
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
