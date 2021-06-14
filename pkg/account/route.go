package account

import (
	"github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
)

type Route struct {
}

func (r *Route) RegisterEndpoint(ws *restful.WebService, dbConnection *gorm.DB) {
	accountHandler := NewAccountHandler(dbConnection)
	ws.Route(ws.GET("/accounts/{accountId}").To(accountHandler.FindByUserID))
	ws.Route(ws.POST("/accounts").To(accountHandler.CreateAccount))
	ws.Route(ws.PUT("/accounts").To(accountHandler.EditAccount))
}
