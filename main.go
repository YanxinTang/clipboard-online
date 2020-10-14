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

func init() {
	execFullPath = os.Args[0]
	execPath = filepath.Dir(execFullPath)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
}

func main() {
	var err error
	// logFileFullPath := execPath + "/" + LogFile
	// f, err := os.OpenFile(logFileFullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	// log.SetOutput(f)

	app, err = NewApplication()
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
