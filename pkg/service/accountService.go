package service

//go:generate mockgen -destination=mock/mockAccountService.go -package=mock github.com/mrth1995/go-mockva/pkg/service AccountService

import (
	"context"
	"fmt"
	"time"

	"github.com/mrth1995/go-mockva/pkg/domain"
	"github.com/mrth1995/go-mockva/pkg/model"
	"github.com/mrth1995/go-mockva/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AccountService defines the interface for account-related business operations.
type AccountService interface {
	// FindByID retrieves an account by its ID.
	FindByID(ctx context.Context, id string) (*domain.Account, error)

	// Register creates a new account.
	Register(ctx context.Context, register *model.AccountRegister) (*domain.Account, error)

	// Edit updates an existing account's information.
	Edit(ctx context.Context, id string, edit *model.AccountEdit) (*domain.Account, error)

	// FindAndLockAccountBalance retrieves an account balance with a pessimistic lock.
	FindAndLockAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error)

	// UpdateBalance updates an account's balance.
	// UpdateBalance updates the account balance within the provided transaction context.
	// Parameters:
	//   - ctx: The request context
	//   - accountBalance: The AccountBalance to update
	//   - tx: The GORM transaction context
	// Returns:
	//   - *domain.AccountBalance: The updated balance
	//   - error: If the operation fails
	UpdateBalance(ctx context.Context, accountBalance *domain.AccountBalance, tx *gorm.DB) (*domain.AccountBalance, error)
}

// AccountServiceImpl implements the AccountService interface.
type AccountServiceImpl struct {
	accountRepository repository.AccountRepository
}

// NewAccountService creates a new instance of AccountService.
func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &AccountServiceImpl{
		accountRepository: accountRepo,
	}
}

func (s *AccountServiceImpl) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	return s.accountRepository.FindByID(ctx, id)
}

func (s *AccountServiceImpl) Register(ctx context.Context, register *model.AccountRegister) (*domain.Account, error) {
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

func (s *AccountServiceImpl) Edit(ctx context.Context, id string, edit *model.AccountEdit) (*domain.Account, error) {
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

func (s *AccountServiceImpl) FindAndLockAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error) {
	return s.accountRepository.FindAndLockAccountBalance(ctx, accountID)
}

// UpdateBalance updates the account balance within the provided transaction context.
// Parameters:
//   - ctx: The request context
//   - accountBalance: The AccountBalance to update
//   - tx: The GORM transaction context
//
// Returns:
//   - *domain.AccountBalance: The updated balance
//   - error: If the operation fails
func (s *AccountServiceImpl) UpdateBalance(ctx context.Context, accountBalance *domain.AccountBalance, tx *gorm.DB) (*domain.AccountBalance, error) {
	return s.accountRepository.UpdateBalance(ctx, accountBalance, tx)
}
