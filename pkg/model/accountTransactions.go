package model

import "time"

type AccountTransactionInfo struct {
	ID                   string    `json:"id"`
	Amount               float64   `json:"amount"`
	AccountSrcID         string    `json:"accountSrcId"`
	AccountSrcName       string    `json:"accountSrcName"`
	AccountDstID         string    `json:"accountDstId"`
	AccountDstName       string    `json:"accountDstName"`
	TransactionTimestamp time.Time `json:"transactionTimestamp"`
}

type AccountFundTransfer struct {
	AccountDstID string  `json:"accountDstId"`
	AccountSrcID string  `json:"accountSrcId"`
	Amount       float64 `json:"amount"`
}
