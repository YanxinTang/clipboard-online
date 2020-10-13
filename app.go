package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lxn/walk"
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
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}
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

func loadConfig() (*Config, error) {
	if isExistFile(ConfigFilePath) {
		return loadConfigFromFile(ConfigFilePath)
	}
	if err := createConfigFile(ConfigFilePath); err != nil {
		return nil, err
	}
	return &DefaultConfig, nil
}

func loadConfigFromFile(path string) (*Config, error) {
	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func createConfigFile(path string) error {
	defaultConfigJson, err := json.MarshalIndent(DefaultConfig, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, []byte(defaultConfigJson), 0644); err != nil {
		return err
	}
	return nil
}

func isExistFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
