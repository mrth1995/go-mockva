package domain

import (
	"time"
)

type Account struct {
	ID        string    `json:"-" gorm:"varchar(32);primaryKey"`
	AccountID string    `json:"accountId" gorm:"varchar(32);not null;unique"`
	Name      string    `json:"name" gorm:"varchar(50);not null"`
	Address   string    `json:"address" gorm:"text"`
	BirthDate time.Time `json:"birthDate" gorm:"not null"`
	Gender    bool      `json:"gender" gorm:"not null"`
}

type AccountBalance struct {
	ID                   string   `json:"-" gorm:"varchar(32);primaryKey"`
	AccountID            string   `json:"accountId" gorm:"varchar(32);not null;unique"`
	Account              *Account `json:"-" gorm:"<-;->:false"`
	Balance              float64  `json:"balance" gorm:"decimal(10,2);not null"`
	AllowNegativeBalance bool     `json:"allowNegativeBalance" gorm:"not null"`
}
