package model

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
	ID        string   `json:"-" gorm:"varchar(32);primaryKey"`
	AccountID string   `json:"accountId" gorm:"varchar(32);not null;unique"`
	Account   *Account `json:"-" gorm:"<-;->:false"`
	Amount    float64  `json:"amount" gorm:"decimal(10,2);not null"`
}

type AccountRegister struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Address              string `json:"address"`
	BirthDate            string `json:"birthDate"`
	Gender               bool   `json:"gender"`
	AllowNegativeBalance bool   `json:"allowNegativeBalance"`
}

type AccountEdit struct {
	Name                 *string `json:"name,omitempty"`
	Address              *string `json:"address,omitempty"`
	BirthDate            *string `json:"birthDate,omitempty"`
	Gender               *bool   `json:"gender,omitempty"`
	AllowNegativeBalance *bool   `json:"allowNegativeBalance,omitempty"`
}
