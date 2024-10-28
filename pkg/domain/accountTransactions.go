package domain

import "time"

type AccountTransaction struct {
	ID                   string          `json:"id" gorm:"varchar(32);primaryKey"`
	TransactionTimestamp time.Time       `json:"transactionTimestamp" gorm:"not null"`
	Amount               float64         `json:"amount" gorm:"not null"`
	AccountSrcId         string          `json:"accountSrcId" gorm:"<-:false;varchar(32);column:accountSrcId"`
	AccountDstId         string          `json:"accountDstId" gorm:"<-:false;varchar(32);column:accountDstId"`
	AccountSrc           *AccountBalance `json:"-" gorm:"<-;->:false"`
	AccountDst           *AccountBalance `json:"-" gorm:"<-;->:false"`
}
