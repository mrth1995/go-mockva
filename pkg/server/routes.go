package server

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/account"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction"
)

func (s *Server) initializeRoutes() {
	s.addHealthCheck()
	s.addRoute(&account.Route{})
	s.addRoute(&accounttransaction.Route{})
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
