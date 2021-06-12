package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AppConfig struct {
	ContextPath string
	Port        int64
}

type DbConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

type Endpoint interface {
	RegisterEndpoint(ws *restful.WebService, db *gorm.DB)
}

type Server struct {
	webService   *restful.WebService
	config       *AppConfig
	listener     net.Listener
	dbConnection *gorm.DB
}

func (s *Server) Initialize(appConfig *AppConfig, dbConfig *DbConfig) {
	s.config = appConfig
	s.initializeDb(dbConfig)
	s.webService = new(restful.WebService)
	s.webService.Path(appConfig.ContextPath)
	s.initializeRoutes()
	restful.Add(s.webService)
}

func (s *Server) Start() {
	var handler http.Handler = restful.DefaultContainer
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return
	}
	logrus.Infof("Server is listening at :%v", s.config.Port)
	err = http.Serve(s.listener, handler)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func (s *Server) Stop() {
	logrus.Infof("Stopping server")
	err := s.listener.Close()
	if err != nil {
		logrus.Errorf("Cannot stop server %v", err)
		return
	}
}

func (s *Server) addRoute(endpoint Endpoint) {
	endpoint.RegisterEndpoint(s.webService, s.dbConnection)
}

func (s *Server) initializeDb(config *DbConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		config.Host, config.Username, config.Password, config.DatabaseName, config.Port)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	s.dbConnection = connection
	for _, entity := range s.RegisterEntities() {
		err := s.dbConnection.AutoMigrate(entity.Entity)
		if err != nil {
			logrus.Error(err)
		}
	}
}
