package accounttransaction

import (
	"time"

	"github.com/mrth1995/go-mockva/pkg/account"
)

type AccountTransaction struct {
	ID                   string           `json:"id" gorm:"varchar(64);primaryKey"`
	TransactionTimestamp time.Time        `json:"transactionTimestamp" gorm:"not null"`
	Amount               float64          `json:"amount" gorm:"not null"`
	AccountSrcId         string           `json:"accountSrcId" gorm:"<-:false;varchar(20);column:accountSrcId"`
	AccountDstId         string           `json:"accountDstId" gorm:"<-:false;varchar(20);column:accountDstId"`
	AccountSrc           *account.Account `json:"-" gorm:"<-;->:false"`
	AccountDst           *account.Account `json:"-" gorm:"<-;->:false"`
}

type AccountFundTransfer struct {
	AccountDstId string  `json:"accountDstId"`
	AccountSrcId string  `json:"accountSrcId"`
	Amount       float64 `json:"amount"`
}
