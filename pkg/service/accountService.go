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
	// FindByID retrieves an account by its unique identifier.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - id: The unique account identifier
	// Returns:
	//   - *domain.Account: The account details if found
	//   - error: If the account is not found or a database error occurs
	FindByID(ctx context.Context, id string) (*domain.Account, error)

	// Register creates a new account in the system.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - register: Account registration details including ID, name, address, birth date, and gender
	// Returns:
	//   - *domain.Account: The newly created account
	//   - error: If account already exists, birth date format is invalid, or database operation fails
	Register(ctx context.Context, register *model.AccountRegister) (*domain.Account, error)

	// Edit updates an existing account's information.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - id: The unique account identifier
	//   - edit: Account fields to update (name, address, gender, birth date) - nil values are ignored
	// Returns:
	//   - *domain.Account: The updated account
	//   - error: If account is not found, birth date format is invalid, or database operation fails
	Edit(ctx context.Context, id string, edit *model.AccountEdit) (*domain.Account, error)

	// FindAndLockAccountBalance retrieves an account balance with a pessimistic lock.
	// This method acquires a database row lock to prevent concurrent modifications during transactions.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - accountID: The unique account identifier
	// Returns:
	//   - *domain.AccountBalance: The account balance with an active lock
	//   - error: If the account is not found or a database error occurs
	FindAndLockAccountBalance(ctx context.Context, accountID string) (*domain.AccountBalance, error)

	// UpdateBalance updates the account balance within the provided transaction context.
	// This method must be called within an active database transaction.
	// Parameters:
	//   - ctx: The request context for cancellation and timeouts
	//   - accountBalance: The AccountBalance entity with the new balance value
	//   - tx: The GORM transaction context
	// Returns:
	//   - *domain.AccountBalance: The updated balance
	//   - error: If the account is not found or a database error occurs
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

// FindByID retrieves an account by its unique identifier.
// Parameters:
//   - ctx: The request context for cancellation and timeouts
//   - id: The unique account identifier
//
// Returns:
//   - *domain.Account: The account details if found
//   - error: If the account is not found or a database error occurs
func (s *AccountServiceImpl) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	return s.accountRepository.FindByID(ctx, id)
}

// Register creates a new account in the system.
// Parameters:
//   - ctx: The request context for cancellation and timeouts
//   - register: Account registration details including ID, name, address, birth date, and gender
//
// Returns:
//   - *domain.Account: The newly created account
//   - error: If account already exists, birth date format is invalid, or database operation fails
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

// Edit updates an existing account's information.
// Parameters:
//   - ctx: The request context for cancellation and timeouts
//   - id: The unique account identifier
//   - edit: Account fields to update (name, address, gender, birth date) - nil values are ignored
//
// Returns:
//   - *domain.Account: The updated account
//   - error: If account is not found, birth date format is invalid, or database operation fails
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

// FindAndLockAccountBalance retrieves an account balance with a pessimistic lock.
// This method acquires a database row lock to prevent concurrent modifications during transactions.
// Parameters:
//   - ctx: The request context for cancellation and timeouts
//   - accountID: The unique account identifier
//
// Returns:
//   - *domain.AccountBalance: The account balance with an active lock
//   - error: If the account is not found or a database error occurs
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
