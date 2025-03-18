package app

import (
	"github.com/ApplyLogic/mqtt-broker/config"
	"github.com/ApplyLogic/mqtt-broker/internal/middleware"
	"github.com/ApplyLogic/mqtt-broker/mqtt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Config *config.Config
	Broker *mqtt.Broker
	Logger *middleware.Logger
}

func (a *App) Initialize(cfg *config.Config) {
	// Initialize Logger
	a.Config = cfg
	a.Logger = &middleware.Logger{}
	a.Logger.Initialize(cfg)
	a.Broker = mqtt.New(a.Config)
}

func (a *App) Start() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	go func() {
		err := a.Broker.Server.Serve()
		if err != nil {
			log.Fatal(err)
		}
		a.Broker.Server.Log.Info("Broker started")
	}()

	<-done
	a.Broker.Server.Log.Warn("caught signal, stopping...")
	_ = a.Broker.Server.Close()
	a.Broker.Server.Log.Info("main.go finished")
}
