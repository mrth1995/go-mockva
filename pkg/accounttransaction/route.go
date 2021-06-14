package accounttransaction

import (
	"github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
)

type Route struct {
}

func (r *Route) RegisterEndpoint(ws *restful.WebService, dbConnection *gorm.DB) {
	h := NewAccountTransactionHandler(dbConnection)
	ws.Route(ws.POST("/accountTransactions/transfer").To(h.Transfer))
}
