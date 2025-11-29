package repository

import "gorm.io/gorm"

// DBTransactionManager defines the interface for managing database transactions.
type DBTransactionManager interface {
	// Transaction executes the given function within a database transaction.
	// If the function returns an error, the transaction is rolled back.
	// Otherwise, the transaction is committed.
	// Parameters:
	//   - fc: The function to execute within the transaction
	// Returns:
	//   - error: If the transaction fails or the function returns an error
	Transaction(fc func(tx *gorm.DB) error) error
}
