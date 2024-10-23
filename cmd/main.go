package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrth1995/go-mockva/pkg/config"
	"github.com/mrth1995/go-mockva/pkg/server"
	"github.com/sirupsen/logrus"
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

	cfg, err := config.ParseConfiguration()
	if err != nil {
		logrus.Fatalf("unable to parse configuration %v", err)
	}

	//setup server
	httpServer := server.Server{}
	httpServer.Initialize(cfg)
	//configure sigint and sigterm
	idleConnectionChan := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownRelease()

		logrus.Infof("OS quit signal received")
		if err = httpServer.Stop(shutdownCtx); err != nil {
			logrus.Errorf("unable to shutdown server: %v", err)
		}
		//handle close database
		logrus.Info("Server stopped")
		close(idleConnectionChan)
	}()
	//serve the server
	if err = httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("unable to server http server: %v", err)
	}
	<-idleConnectionChan
}
