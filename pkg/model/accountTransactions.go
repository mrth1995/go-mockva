package model

import "time"

type AccountTransactionInfo struct {
	ID                   string    `json:"id"`
	Amount               float64   `json:"amount"`
	AccountSrcId         string    `json:"accountSrcId"`
	AccountSrcName       string    `json:"accountSrcName"`
	AccountDstId         string    `json:"accountDstId"`
	AccountDstName       string    `json:"accountDstName"`
	TransactionTimestamp time.Time `json:"transactionTimestamp"`
}

type AccountFundTransfer struct {
	AccountDstId string  `json:"accountDstId"`
	AccountSrcId string  `json:"accountSrcId"`
	Amount       float64 `json:"amount"`
}
