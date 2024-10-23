package service

import (
	"fmt"
	"time"

	"github.com/mrth1995/go-mockva/pkg/account/model"
	"github.com/mrth1995/go-mockva/pkg/account/repository"
	"github.com/sirupsen/logrus"
)

type AccountServiceImpl struct {
	repository repository.AccountRepository
}

func (s *AccountServiceImpl) Register(register *model.AccountRegister) (*model.Account, error) {
	existingAccount, notFound := s.repository.FindById(register.ID)
	if existingAccount != nil && notFound == nil {
		return nil, fmt.Errorf("account %v already exist", register.ID)
	}

	birthDate, err := time.Parse(time.DateOnly, register.BirthDate)
	if err != nil {
		logrus.Errorf("Invalid date format %v", register.BirthDate)
		return nil, err
	}

	newAccount := &model.Account{
		ID:        register.ID,
		Name:      register.Name,
		Address:   register.Address,
		BirthDate: birthDate,
		Gender:    register.Gender,
	}
	accountAlreadyExist := s.repository.Save(newAccount)
	if accountAlreadyExist != nil {
		return nil, accountAlreadyExist
	}
	return newAccount, nil
}

func (s *AccountServiceImpl) Edit(id string, edit *model.AccountEdit) (*model.Account, error) {
	existingAccount, err := s.repository.FindById(id)
	if err != nil {
		return nil, err
	}
	if edit.Name != nil {
		existingAccount.Name = *edit.Name
	}
	if edit.Address != nil {
		existingAccount.Address = *edit.Address
	}
	if edit.Gender != nil {
		existingAccount.Gender = *edit.Gender
	}
	if edit.BirthDate != nil && *edit.BirthDate != "" {
		birthDate, err := time.Parse(time.DateOnly, *edit.BirthDate)
		if err != nil {
			logrus.Errorf("Invalid date format %v", *edit.BirthDate)
			return nil, err
		}
		existingAccount.BirthDate = birthDate
	}
	_, err = s.repository.Update(existingAccount)
	if err != nil {
		return nil, err
	}
	return existingAccount, nil
}
