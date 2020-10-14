package main

import (
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lxn/walk"
	log "github.com/sirupsen/logrus"
)

func router() *httprouter.Router {
	router := httprouter.New()
	router.GET("/clipboard", getHandler)
	router.POST("/clipboard", setHandler)
	return router
}

func getHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requestLogger := log.WithFields(log.Fields{"request_id": rand.Int(), "user_ip": r.RemoteAddr})
	str, err := walk.Clipboard().Text()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "")
		requestLogger.WithError(err).Warn("failed to get clipboard")
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, str)
	requestLogger.Info("get clipboard text")
}

func setHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requestLogger := log.WithFields(log.Fields{"request_id": rand.Int(), "user_ip": r.RemoteAddr})

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		requestLogger.WithError(err).Warn("failed to read request body")
		return
	}
	defer r.Body.Close()
	bodystr := string(body)

	if err := walk.Clipboard().SetText(bodystr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		requestLogger.WithError(err).Warn("failed to set clipboard")
		return
	}

	notify := bodystr
	if notify == "" {
		notify = "粘贴内容为空"
	}

	if err := app.ni.ShowInfo("粘贴自我的设备", notify); err != nil {
		requestLogger.WithError(err).WithField("notify", notify).Warn("failed to send notification")
	}

	w.WriteHeader(http.StatusOK)
	requestLogger.WithField("text", bodystr).Info("set clipboard text")
	return
}
