package server

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type Config struct {
	ContextPath string
	Port int64
	IPAddr string
}

type Server struct {
	container *restful.Container
	config *Config
	listener net.Listener
}

func (s *Server) Serve() {
	var handler http.Handler = s.container
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.IPAddr, s.config.Port))
	if err != nil {
		return
	}
	logrus.Infof("Server is listening at %v:%v", s.config.IPAddr, s.config.Port)
	err = http.Serve(s.listener, handler)
	if err != nil {
		return
	}
}

func (s *Server) addHealthCheck(ws *restful.WebService) {
	ws.Route(ws.GET("/healthz").To(func(request *restful.Request, response *restful.Response) {
		content := map[string]interface{}{
			"healthy": true,
		}
		err := response.WriteAsJson(content)
		if err != nil {
			return
		}
	}))
}

func (s *Server) Stop() {
	logrus.Infof("Stopping server")
	err := s.listener.Close()
	if err != nil {
		logrus.Error(err)
		return
	}
}

func New(cfg *Config) *Server {
	s := &Server{
		container: restful.DefaultContainer,
		config: cfg,
	}
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{},
		AllowedHeaders: []string{
			"Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Authorization",
			"Content-Type", "Accept",
		},
		CookiesAllowed: true,
		Container:      s.container,
	}
	s.container.Filter(cors.Filter)
	rootWebService := new(restful.WebService)
	s.addHealthCheck(rootWebService)
	s.container.Add(rootWebService)
	return s
}
