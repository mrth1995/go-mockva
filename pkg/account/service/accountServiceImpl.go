package service

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mrth1995/go-mockva/pkg/account/model"
	"github.com/mrth1995/go-mockva/pkg/account/repository"
)

type AccountServiceImpl struct {
	repository repository.AccountRepository
}

func (s *AccountServiceImpl) Register(register *model.AccountRegister) (*model.Account, error) {
	existingAccount, notFound := s.repository.FindById(register.ID)
	if existingAccount != nil && notFound == nil {
		return nil, fmt.Errorf("account %v already exist", register.ID)
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

	newAccount := &model.Account{
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

func (s *AccountServiceImpl) Edit(edit *model.AccountEdit) (*model.Account, error) {
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
