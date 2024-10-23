package route

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/account/handler"
	"gorm.io/gorm"
)

type AccountRoute struct {
}

func (r *AccountRoute) RegisterEndpoint(ws *restful.WebService, dbConnection *gorm.DB) {
	accountHandler := handler.NewAccountHandler(dbConnection)
	ws.Route(ws.GET("/accounts/{accountId}").To(accountHandler.FindByUserID))
	ws.Route(ws.POST("/accounts").To(accountHandler.CreateAccount))
	ws.Route(ws.PUT("/accounts/{accountId}").To(accountHandler.EditAccount))
}
