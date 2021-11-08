package main

import (
	"path"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lxn/walk"
)

type Application struct {
	config *Config
	*walk.MainWindow
	ni *walk.NotifyIcon
	wg sync.WaitGroup
}

func (app *Application) RunHTTPServer() {
	app.wg.Add(1)
	go func() {
		engin := gin.New()
		setupRoute(engin)
		if err := engin.Run(":" + app.config.Port); err != nil {
			app.ni.ShowError("HTTP Server 启动失败", "您的应用可能不能正常运行")
			app.Synchronize(func() {
				walk.App().Exit(1)
			})
			log.WithError(err).Error("failed to start http server")
			return
		}
	}()
}

func (app *Application) StopHTTPServer() {
	app.wg.Done()
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

func (app *Application) GetTempFilePath(filename string) string {
	if !filepath.IsAbs(app.config.TempDir) {
		// temp files path in exec path but not pwd
		tempAbsPath := path.Join(execPath, app.config.TempDir)
		return filepath.Join(tempAbsPath, filename)
	}
	return filepath.Join(app.config.TempDir, filename)
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

	return app, nil
}
