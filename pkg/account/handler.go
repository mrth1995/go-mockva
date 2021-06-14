package account

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/server/responseWriter"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	AccountService    Service
	AccountRepository Repository
}

func (h *Handler) FindByUserID(request *restful.Request, response *restful.Response) {
	accountId := request.PathParameter("accountId")
	existingAccount, err := h.AccountRepository.FindById(accountId)
	if existingAccount == nil && err != nil {
		logrus.Infof("Account %v not found", accountId)
		responseWriter.WriteBadRequest(err, response)
		return
	}
	responseWriter.WriteOK(existingAccount, response)
}

func (h *Handler) CreateAccount(request *restful.Request, response *restful.Response) {
	var accountRegister Register
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

func (h *Handler) EditAccount(request *restful.Request, response *restful.Response) {
	var accountEdit Edit
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

func NewAccountHandler(dbConnection *gorm.DB) *Handler {
	return &Handler{
		AccountRepository: NewRepository(dbConnection),
		AccountService:    NewService(NewRepository(dbConnection)),
	}
}
