// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/YanxinTang/clipboard-online/action"
	"github.com/lxn/walk"
)

var app *Application

func main() {
	var err error
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

	exitAction, err := action.NewExitAction()
	if err != nil {
		log.Fatal(err)
	}
	if err := app.AddActions(exitAction); err != nil {
		log.Fatal(err)
	}

	if err := app.ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	app.RunHTTPServer()
	app.Run()
}
