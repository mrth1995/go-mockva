package handler

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/account/repository"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/model"
	trxRepo "github.com/mrth1995/go-mockva/pkg/accounttransaction/repository"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/service"
	"github.com/mrth1995/go-mockva/pkg/server/responseWriter"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountTransactionHandler struct {
	AccountTransactionService service.AccountTransactionService
}

func NewAccountTransactionHandler(dbConn *gorm.DB) *AccountTransactionHandler {
	accountTrxService := &service.AccountTransactionServiceImpl{
		AccountRepository:    repository.NewAccountRepository(dbConn),
		AccountTrxRepository: trxRepo.NewAccountTrxRepository(dbConn),
	}
	return &AccountTransactionHandler{AccountTransactionService: accountTrxService}
}

func (h *AccountTransactionHandler) Transfer(request *restful.Request, response *restful.Response) {
	var param model.AccountFundTransfer
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
