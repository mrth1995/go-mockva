package server

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/account/route"
	trxRoute "github.com/mrth1995/go-mockva/pkg/accounttransaction/route"
)

func (s *Server) initializeRoutes() {
	s.addHealthCheck()
	s.addRoute(&route.AccountRoute{})
	s.addRoute(&trxRoute.AccountTransactionRoute{})
}

func (s *Server) addHealthCheck() {
	s.webService.Route(s.webService.GET("/healthz").To(func(request *restful.Request, response *restful.Response) {
		content := map[string]interface{}{
			"healthy": true,
		}
		err := response.WriteAsJson(content)
		if err != nil {
			return
		}
	}))
}
