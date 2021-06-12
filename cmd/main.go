package main

import (
	"github.com/mrth1995/go-mockva/pkg/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const millisecondTimeFormat = "2006-01-02T15:04:05.999Z07:00"

func main() {
	//setup logrus
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: millisecondTimeFormat,
	})

	//setup server
	appConfig := &server.AppConfig{
		ContextPath: "/mockva",
		Port:        8080,
	}
	dbConfig := &server.DbConfig{
		Host:         "localhost",
		Port:         5432,
		Username:     "mrth1995",
		Password:     "karuiongaku123",
		DatabaseName: "go-mockva",
	}
	httpServer := server.Server{}
	httpServer.Initialize(appConfig, dbConfig)
	//configure sigint and sigterm
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		logrus.Infof("OS quit signal received")
		httpServer.Stop()
		//handle close database
		logrus.Info("Server stopped")
	}()
	//serve the server
	httpServer.Start()
}
