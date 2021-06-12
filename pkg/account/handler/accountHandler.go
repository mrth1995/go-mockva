package handler

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/account/model"
	"github.com/mrth1995/go-mockva/pkg/account/repository"
	"github.com/mrth1995/go-mockva/pkg/account/service"
	"github.com/mrth1995/go-mockva/pkg/server/responseWriter"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountHandler struct {
	AccountService    service.AccountService
	AccountRepository repository.AccountRepository
}

func (h *AccountHandler) FindByUserID(request *restful.Request, response *restful.Response) {
	accountId := request.PathParameter("accountId")
	existingAccount, err := h.AccountRepository.FindById(accountId)
	if existingAccount == nil && err != nil {
		logrus.Infof("Account %v not found", accountId)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	responseWriter.WriteOK(existingAccount, response)
}

func (h *AccountHandler) CreateAccount(request *restful.Request, response *restful.Response) {
	var accountRegister model.AccountRegister
	err := request.ReadEntity(&accountRegister)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	newAccount, err := h.AccountService.Register(&accountRegister)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	responseWriter.WriteOK(newAccount, response)
}

func (h *AccountHandler) EditAccount(request *restful.Request, response *restful.Response) {
	var accountEdit model.AccountEdit
	err := request.ReadEntity(accountEdit)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	account, err := h.AccountService.Edit(&accountEdit)
	if err != nil {
		logrus.Error(err)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	responseWriter.WriteOK(account, response)
}

func NewAccountHandler(dbConnection *gorm.DB) *AccountHandler {
	return &AccountHandler{
		AccountRepository: repository.NewAccountRepository(dbConnection),
		AccountService:    service.NewAccountService(repository.NewAccountRepository(dbConnection)),
	}
}
