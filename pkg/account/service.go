package account

import (
	"time"

	"github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Register(register *Register) (*Account, error)
	Edit(edit *Edit) (*Account, error)
}

type serviceImpl struct {
	repository Repository
}

func NewService(accountRepository Repository) Service {
	return &serviceImpl{repository: accountRepository}
}

func (s *serviceImpl) Register(register *Register) (*Account, error) {
	existingAccount, notFound := s.repository.FindById(register.ID)
	if existingAccount != nil && notFound == nil {
		return nil, errors.NewAccountAlreadyExist(register.ID)
	}
	var birthDate *time.Time = nil
	if register.BirthDate != "" {
		parse, err := time.Parse("2006-01-02", register.BirthDate)
		birthDate = &parse
		if err != nil {
			logrus.Errorf("Invalid date format %v", register.BirthDate)
			return nil, err
		}
	}

	newAccount := &Account{
		ID:                   register.ID,
		Name:                 register.Name,
		Address:              register.Address,
		BirthDate:            birthDate,
		Gender:               register.Gender,
		AllowNegativeBalance: register.AllowNegativeBalance,
	}
	accountAlreadyExist := s.repository.Save(newAccount)
	if accountAlreadyExist != nil {
		return nil, accountAlreadyExist
	}
	return newAccount, nil
}

func (s *serviceImpl) Edit(edit *Edit) (*Account, error) {
	existingAccount, err := s.repository.FindById(edit.ID)
	if err != nil {
		return nil, err
	}
	existingAccount.Name = edit.Name
	existingAccount.AllowNegativeBalance = edit.AllowNegativeBalance
	existingAccount.Address = edit.Address
	existingAccount.Gender = edit.Gender
	if edit.BirthDate != "" {
		updatedBirthDate, err := time.Parse("2006-01-02", edit.BirthDate)
		if err != nil {
			return nil, err
		}
		existingAccount.BirthDate = &updatedBirthDate
	}
	_, err = s.repository.Update(existingAccount)
	if err != nil {
		return nil, err
	}
	return existingAccount, nil
}
