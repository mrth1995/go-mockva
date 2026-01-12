package controller

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/model"
	"github.com/mrth1995/go-mockva/pkg/server/responseWriter"
	"github.com/mrth1995/go-mockva/pkg/service"
	"github.com/sirupsen/logrus"
)

type AccountController struct {
	AccountService service.AccountService
}

func NewAccountController(accountService service.AccountService) *AccountController {
	return &AccountController{
		AccountService: accountService,
	}
}

func (accountController *AccountController) FindByUserID(request *restful.Request, response *restful.Response) {
	ctx := request.Request.Context()

	accountID := request.PathParameter("accountId")
	existingAccount, err := accountController.AccountService.FindByID(ctx, accountID)
	if existingAccount == nil && err != nil {
		logrus.Infof("Account %v not found", accountID)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	responseWriter.WriteOK(existingAccount, response)
}

func (accountController *AccountController) CreateAccount(request *restful.Request, response *restful.Response) {
	ctx := request.Request.Context()

	var accountRegister model.AccountRegister
	err := request.ReadEntity(&accountRegister)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	newAccount, err := accountController.AccountService.Register(ctx, &accountRegister)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	responseWriter.WriteOK(newAccount, response)
}

func (accountController *AccountController) EditAccount(request *restful.Request, response *restful.Response) {
	ctx := request.Request.Context()

	accountID := request.PathParameter("accountId")

	var accountEdit model.AccountEdit
	err := request.ReadEntity(accountEdit)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	account, err := accountController.AccountService.Edit(ctx, accountID, &accountEdit)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	responseWriter.WriteOK(account, response)
}
