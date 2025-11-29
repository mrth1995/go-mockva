package postgresql

import (
	"github.com/mrth1995/go-mockva/pkg/repository"
	"gorm.io/gorm"
)

// GormTransactionManager implements DBTransactionManager using GORM.
type GormTransactionManager struct {
	db *gorm.DB
}

// NewGormTransactionManager creates a new GORM transaction manager.
func NewGormTransactionManager(db *gorm.DB) repository.DBTransactionManager {
	return &GormTransactionManager{db: db}
}

// Transaction executes the given function within a GORM transaction.
func (m *GormTransactionManager) Transaction(fc func(tx *gorm.DB) error) error {
	return m.db.Transaction(fc)
}
