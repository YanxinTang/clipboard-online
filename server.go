package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lxn/walk"
)

func router() *httprouter.Router {
	router := httprouter.New()
	router.GET("/clipboard", getHandler)
	router.POST("/clipboard", setHandler)
	return router
}

func getHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	str, err := walk.Clipboard().Text()
	if err != nil {
		io.WriteString(w, "")
		log.Printf("[ERROR]: %s %s", r.Method, err)
	}
	log.Printf("[INFO]: %s", r.Method)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, str)
}

func setHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	bodystr := string(body)
	if err != nil {
		log.Printf("[ERROR]: %s %s", r.Method, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := walk.Clipboard().SetText(bodystr); err != nil {
		log.Printf("[ERROR]: %s %s", r.Method, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("[INFO]: %s", r.Method)

	notify := bodystr
	if notify == "" {
		notify = "粘贴内容为空"
	}
	if err := app.ni.ShowInfo("粘贴自我的设备", notify); err != nil {
		log.Printf("[ERROR]: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	return
}
