package model

import "time"

type AccountInfo struct {
	ID        string    `json:"-"`
	AccountID string    `json:"accountId"`
	Amount    float64   `json:"amount"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	BirthDate time.Time `json:"birthDate"`
	Gender    bool      `json:"gender"`
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
