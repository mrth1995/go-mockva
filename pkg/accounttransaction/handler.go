package accounttransaction

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/account"
	"github.com/mrth1995/go-mockva/pkg/server/responseWriter"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	AccountTransactionService Service
}

func NewAccountTransactionHandler(dbConn *gorm.DB) *Handler {
	accountTrxService := &serviceImpl{
		AccountRepository:    account.NewRepository(dbConn),
		AccountTrxRepository: NewRepository(dbConn),
	}
	return &Handler{AccountTransactionService: accountTrxService}
}

func (h *Handler) Transfer(request *restful.Request, response *restful.Response) {
	var param AccountFundTransfer
	err := request.ReadEntity(&param)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	trx, err := h.AccountTransactionService.Transfer(&param)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	responseWriter.WriteOK(trx, response)
}
