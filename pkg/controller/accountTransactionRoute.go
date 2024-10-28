package controller

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	endpointError "github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/mrth1995/go-mockva/pkg/model"
	"gorm.io/gorm"
)

func (accountTransactionHandler *AccountTransactionController) RegisterEndpoint(ws *restful.WebService, dbConnection *gorm.DB) {
	tags := []string{"Account Transactions"}
	ws.Route(
		ws.POST("/accountTransactions/transfer").
			To(accountTransactionHandler.Transfer).
			Consumes(restful.MIME_JSON).
			Produces(restful.MIME_JSON).
			Reads(model.AccountFundTransfer{}).
			Returns(http.StatusOK, "Transaction success", model.AccountTransactionInfo{}).
			Returns(http.StatusBadRequest, "Validation error", endpointError.EndpointError{}).
			Returns(http.StatusInternalServerError, "Internal server error", endpointError.EndpointError{}).
			Metadata(restfulspec.KeyOpenAPITags, tags))
}
