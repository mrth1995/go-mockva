package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/mrth1995/go-mockva/pkg/model"
	"github.com/mrth1995/go-mockva/pkg/repository"
	"github.com/sirupsen/logrus"
)

type AccountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepository: accountRepo,
	}
}

func (s *AccountService) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	// todo: implement me
	return nil, errors.New("unsupported operation")
}

func (s *AccountService) Register(ctx context.Context, register *model.AccountRegister) (*domain.Account, error) {
	existingAccount, notFound := s.accountRepository.FindByID(ctx, register.ID)
	if existingAccount != nil && notFound == nil {
		return nil, fmt.Errorf("account %v already exist", register.ID)
	}

	birthDate, err := time.Parse(time.DateOnly, register.BirthDate)
	if err != nil {
		logrus.Errorf("Invalid date format %v", register.BirthDate)
		return nil, err
	}

	newAccount := &domain.Account{
		ID:        register.ID,
		Name:      register.Name,
		Address:   register.Address,
		BirthDate: birthDate,
		Gender:    register.Gender,
	}
	accountAlreadyExist := s.accountRepository.Save(ctx, newAccount)
	if accountAlreadyExist != nil {
		return nil, accountAlreadyExist
	}
	return newAccount, nil
}

func (s *AccountService) Edit(ctx context.Context, id string, edit *model.AccountEdit) (*domain.Account, error) {
	existingAccount, err := s.accountRepository.FindByID(ctx, id)
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
	_, err = s.accountRepository.Update(ctx, existingAccount)
	if err != nil {
		return nil, err
	}
	return existingAccount, nil
}

func (s *AccountService) FindAndLockAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error) {
	// todo: implement me
	return nil, errors.New("unsupported operation")
}

func (s *AccountService) UpdateBalance(ctx context.Context, accountBalance *domain.AccountBalance) (*domain.AccountBalance, error) {
	// todo: implement me
	return nil, errors.New("unsupported operation")
}
