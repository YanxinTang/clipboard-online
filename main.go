// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"

	"github.com/YanxinTang/clipboard-online/action"
	"github.com/lxn/walk"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var app *Application

var execPath string
var execFullPath string
var config *Config
var mode string = "debug"

func init() {
	execFullPath = os.Args[0]
	execPath = filepath.Dir(execFullPath)

	var err error
	configFilePath := filepath.Join(execPath, ConfigFile)
	config, err = loadConfig(configFilePath)
	if err != nil {
		log.WithError(err).Warn("failed to load config")
	}
	log.SetLevel(config.LogLevel)

	if mode == "debug" {
		log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	} else {
		logFilePath := filepath.Join(execPath, LogFile)
		f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.WithError(err).Fatal("failed to open log file")
		}
		log.SetOutput(f)
	}
}

func main() {
	var err error
	app, err = NewApplication(config)
	if err != nil {
		log.WithError(err).Fatal("failed to create applicaton")
	}
	defer app.BeforeExit()

	icon, err := walk.NewIconFromResourceId(3)
	if err != nil {
		log.WithError(err).Fatal("failed to get icon")
	}

	if err := app.ni.SetIcon(icon); err != nil {
		log.WithError(err).Fatal("failed to set icon")
	}

	if err := app.ni.SetToolTip("clipboard-online"); err != nil {
		log.WithError(err).Fatal("failed to set tooltip")
	}

	autoRunAction, err := action.NewAutoRunAction()
	if err != nil {
		log.WithError(err).Fatal("failed to create AutoRunAction")
	}
	exitAction, err := action.NewExitAction()
	if err != nil {
		log.WithError(err).Fatal("failed to create ExitAction")
	}
	if err := app.AddActions(autoRunAction, exitAction); err != nil {
		log.WithError(err).Fatal("failed to add action")
	}

	if err := app.ni.SetVisible(true); err != nil {
		log.WithError(err).Fatal("failed to set notify visible")
	}

	log.Debug("start http server")
	app.RunHTTPServer()
	log.Debug("start app")
	app.Run()
}
