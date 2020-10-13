// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/YanxinTang/clipboard-online/action"
	"github.com/lxn/walk"
)

var app *Application

var execPath string
var execFullPath string

func init() {
	execFullPath = os.Args[0]
	execPath = filepath.Dir(execFullPath)
}

func main() {
	var err error

	logFileFullPath := execPath + "/" + LogFile
	f, err := os.OpenFile(logFileFullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	app, err = NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	defer app.BeforeExit()

	icon, err := walk.NewIconFromResourceId(3)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}

	if err := app.ni.SetToolTip("clipboard-online"); err != nil {
		log.Fatal(err)
	}

	autoRunAction, err := action.NewAutoRunAction()
	if err != nil {
		log.Fatal(err)
	}
	exitAction, err := action.NewExitAction()
	if err != nil {
		log.Fatal(err)
	}
	if err := app.AddActions(autoRunAction, exitAction); err != nil {
		log.Fatal(err)
	}

	if err := app.ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}
	app.RunHTTPServer()
	app.Run()
}
