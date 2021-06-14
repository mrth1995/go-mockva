package account

import (
	"time"
)

type Account struct {
	ID                   string     `json:"id" gorm:"varchar(64);primaryKey"`
	Name                 string     `json:"name" gorm:"varchar(50);not null"`
	Address              string     `json:"address" gorm:"text"`
	BirthDate            *time.Time `json:"birthDate"`
	Gender               bool       `json:"gender" gorm:"not null"`
	Balance              float64    `json:"balance" gorm:"not null"`
	AllowNegativeBalance bool       `json:"allowNegativeBalance" gorm:"not null"`
}

type Register struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Address              string `json:"address"`
	BirthDate            string `json:"birthDate"`
	Gender               bool   `json:"gender"`
	AllowNegativeBalance bool   `json:"allowNegativeBalance"`
}

type Edit struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Address              string `json:"address"`
	BirthDate            string `json:"birthDate"`
	Gender               bool   `json:"gender"`
	AllowNegativeBalance bool   `json:"allowNegativeBalance"`
}

func (a *Account) GetTableName() string {
	return "account"
}
