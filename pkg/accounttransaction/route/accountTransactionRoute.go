package route

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/handler"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/model"
	endpointError "github.com/mrth1995/go-mockva/pkg/errors"
	"gorm.io/gorm"
)

type AccountTransactionRoute struct {
}

func (r *AccountTransactionRoute) RegisterEndpoint(ws *restful.WebService, dbConnection *gorm.DB) {
	tags := []string{"Account Transactions"}
	h := handler.NewAccountTransactionHandler(dbConnection)
	ws.Route(
		ws.POST("/accountTransactions/transfer").
			To(h.Transfer).
			Consumes(restful.MIME_JSON).
			Produces(restful.MIME_JSON).
			Reads(model.AccountFundTransfer{}).
			Returns(http.StatusOK, "Transaction success", model.AccountTransaction{}).
			Returns(http.StatusBadRequest, "Validation error", endpointError.EndpointError{}).
			Returns(http.StatusInternalServerError, "Internal server error", endpointError.EndpointError{}).
			Metadata(restfulspec.KeyOpenAPITags, tags))
}
