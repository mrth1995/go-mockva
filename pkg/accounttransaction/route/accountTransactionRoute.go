package route

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction/handler"
	"gorm.io/gorm"
)

type AccountTransactionRoute struct {
}

func (r *AccountTransactionRoute) RegisterEndpoint(ws *restful.WebService, dbConnection *gorm.DB) {
	h := handler.NewAccountTransactionHandler(dbConnection)
	ws.Route(ws.POST("/accountTransactions/transfer").To(h.Transfer))
}
