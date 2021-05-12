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
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: millisecondTimeFormat,
		FullTimestamp: true,
	})
	serverConfig := &server.Config{
		ContextPath: "/mockva",
		Port:        8080,
	}
	httpServer := server.New(serverConfig)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		logrus.Infof("OS quit signal received")
		httpServer.Stop()
		logrus.Info("Server stopped")
	}()
	httpServer.Serve()
}
