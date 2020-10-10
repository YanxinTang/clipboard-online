package main

import (
	"log"
	"net/http"

	"github.com/lxn/walk"
)

type Application struct {
	*walk.MainWindow
	ni         *walk.NotifyIcon
	serverChan chan string
}

func (app *Application) RunHTTPServer() {
	go func() {
		router := router()
		if err := http.ListenAndServe(":8000", router); err != nil {
			log.Print(err)
			app.ni.ShowError("HTTP Server 启动失败", "您的应用可能不能正常运行")
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

func NewApplication() (*Application, error) {
	app := new(Application)
	var err error
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
