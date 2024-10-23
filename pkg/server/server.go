package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/config"
	"github.com/mrth1995/go-mockva/pkg/migration"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	contextPath = "/mockva"
)

type Endpoint interface {
	RegisterEndpoint(ws *restful.WebService, db *gorm.DB)
}

type Server struct {
	webService   *restful.WebService
	cfg          *config.Config
	httpServer   *http.Server
	dbConnection *gorm.DB
}

func (s *Server) Initialize(cfg *config.Config) {
	s.cfg = cfg
	s.initializeDb()
	s.migrateDBSchema()
	s.webService = new(restful.WebService)
	s.webService.Path(contextPath)
	s.initializeRoutes()
	restful.Add(s.webService)
}

func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Handler: restful.DefaultContainer,
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.cfg.Port),
	}
	logrus.Infof("Server is listening at :%v", s.cfg.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	logrus.Infof("Stopping server")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) addRoute(endpoint Endpoint) {
	endpoint.RegisterEndpoint(s.webService, s.dbConnection)
}

func (s *Server) initializeDb() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		s.cfg.PostgresHost, s.cfg.PostgresUsername, s.cfg.PostgresPassword, s.cfg.DBName, s.cfg.PostgresPort)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	s.dbConnection = connection
}

func (s *Server) migrateDBSchema() {
	DB, _ := s.dbConnection.DB()
	migration, err := migration.NewMigration(DB, s.cfg.DBName, s.cfg.SQLFilePath+"/postgresql")
	if err != nil {
		logrus.Fatal(err)
	}
	if err = migration.Up(); err != nil {
		logrus.Fatal(err)
	}
}
