package server

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/mrth1995/go-mockva/pkg/controller"
	"github.com/mrth1995/go-mockva/pkg/repository/postgresql"
	"github.com/mrth1995/go-mockva/pkg/service"
)

func (s *Server) initializeRoutes() {
	ws := new(restful.WebService)
	ws.Path(contextPath)

	accountRepository := postgresql.NewAccountRepository(s.dbConnection)
	accountTrxRepository := postgresql.NewAccountTrxRepository(s.dbConnection)

	accountService := service.NewAccountService(accountRepository)
	accountTrxService := service.NewAccountTrxService(accountService, accountTrxRepository)

	accountController := controller.NewAccountController(accountService)
	accountTrxController := controller.NewAccountTransactionController(accountTrxService)

	s.addRoute(ws, accountController)
	s.addRoute(ws, accountTrxController)
	restful.Add(ws)
	s.addSwaggerDocs()
}

func (s *Server) addSwaggerDocs() {
	webServices := restful.DefaultContainer.RegisteredWebServices()
	swaggerConfig := restfulspec.Config{
		WebServices: webServices,
		APIPath:     contextPath + "/apidocs/api.json",
		PostBuildSwaggerObjectHandler: func(s *spec.Swagger) {
			s.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "mockva",
					Description: "Mock Virtual Account",
					Version:     "1.0.0",
					Contact: &spec.ContactInfo{
						ContactInfoProps: spec.ContactInfoProps{
							Name:  "M Ridwan Taufik H",
							Email: "mr.taufikhidayat.1995@gmail.com",
							URL:   "https://www.linkedin.com/in/m-ridwan-taufik-hidayat-775765138/",
						},
					},
				},
			}
		},
	}

	http.Handle(contextPath+"/apidocs/", http.StripPrefix(contextPath+"/apidocs/", http.FileServer(http.Dir(s.cfg.SwaggerFilePath))))

	service := restfulspec.NewOpenAPIService(swaggerConfig)

	restful.Add(service)
}
