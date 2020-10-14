package main

import (
	"net/http"

	"github.com/lxn/walk"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	config *Config
	*walk.MainWindow
	ni         *walk.NotifyIcon
	serverChan chan string
}

func (app *Application) RunHTTPServer() {
	go func() {
		router := router()
		if err := http.ListenAndServe(":"+app.config.Port, router); err != nil {
			app.ni.ShowError("HTTP Server 启动失败", "您的应用可能不能正常运行")
			log.WithError(err).Error("failed to start http server")
			return
		}
		for range app.serverChan {
		}
	}()
}

func (app *Application) StopHTTPServer() {
	close(app.serverChan)
}

func (app *Application) BeforeExit() {
	app.StopHTTPServer()
	app.ni.Dispose()
}

func (app *Application) AddActions(actions ...*walk.Action) error {
	for _, action := range actions {
		if err := app.ni.ContextMenu().Actions().Add(action); err != nil {
			return err
		}
	}
	return nil
}

func NewApplication(config *Config) (*Application, error) {
	app := new(Application)
	var err error
	app.config = config
	app.MainWindow, err = walk.NewMainWindow()
	if err != nil {
		return nil, err
	}

	app.ni, err = walk.NewNotifyIcon(app.MainWindow)
	if err != nil {
		return nil, err
	}

	app.serverChan = make(chan string)
	return app, nil
}
