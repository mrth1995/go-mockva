// Package controller is package for http controller
package controller

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/model"
	"github.com/mrth1995/go-mockva/pkg/server/responseWriter"
	"github.com/mrth1995/go-mockva/pkg/service"
	"github.com/sirupsen/logrus"
)

type AccountTransactionController struct {
	AccountTransactionService *service.AccountTransactionService
}

func NewAccountTransactionController(accountTransactionService *service.AccountTransactionService) *AccountTransactionController {
	return &AccountTransactionController{AccountTransactionService: accountTransactionService}
}

// Transfer move account balance from one account to another account
func (accountTransactionController *AccountTransactionController) Transfer(request *restful.Request, response *restful.Response) {
	ctx := request.Request.Context()

	var param model.AccountFundTransfer
	err := request.ReadEntity(&param)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	trx, err := accountTransactionController.AccountTransactionService.Transfer(ctx, &param)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	responseWriter.WriteOK(trx, response)
}
