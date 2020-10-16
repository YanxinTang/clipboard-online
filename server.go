package main

import (
	"bufio"
	"github.com/YanxinTang/clipboard-online/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/lxn/walk"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

func router() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", getHandler)
	router.POST("/", setHandler)
	router.NotFound = http.HandlerFunc(notFoundHandler)
	return router
}

func getHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requestLogger := log.WithFields(log.Fields{"request_id": rand.Int(), "user_ip": r.RemoteAddr})
	str, err := walk.Clipboard().Text()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "")
		requestLogger.WithError(err).Warn("failed to get clipboard")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, str)
	requestLogger.Info("get clipboard text")
}

const (
	typeText  = "text"
	typeFile  = "file"
	typeMedia = "media"
)

func setHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requestLogger := log.WithFields(log.Fields{"request_id": rand.Int(), "user_ip": r.RemoteAddr})
	rd := bufio.NewReader(r.Body)

	var (
		version     string
		contentType string
		notify      string
		filename    string
	)

	q := r.URL.Query()
	version = q.Get("version")
	contentType = q.Get("type")
	filename = q.Get("filename")

	if len(filename) == 0 {
		filename = utils.RandStringBytes(16)
	}

	cleanTempFile()

	if version == "" || contentType == typeText {
		text, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			requestLogger.WithError(err).Warn("failed to read request body")
			return
		}

		if len(text) == 0 {
			notify = "粘贴内容为空"
		} else {
			notify = string(text)
		}

		if err := walk.Clipboard().SetText(string(text)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			requestLogger.WithError(err).Warn("failed to set clipboard")
			return
		}

		requestLogger.WithField("text", string(text)).Info("set clipboard text")
	} else if contentType == typeFile || contentType == typeMedia {
		path := getTempFilePath(filename)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			requestLogger.WithError(err).Warn("failed create temporary file")
			return
		}
		defer file.Close()

		_, _ = io.Copy(file, rd)

		if contentType == typeMedia {
			notify = "[图片媒体] 已复制到剪贴板"
		} else {
			notify = "[文件] 已复制到剪贴板"
		}

		setLastFilename(filename)

		if err := utils.Clipboard().SetFile(path); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			requestLogger.WithError(err).Warn("failed to set clipboard")
			return
		}

		requestLogger.WithField("file", path).Info("set clipboard file")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		requestLogger.Warn("unsupported content type")
		return
	}

	if err := app.ni.ShowInfo("粘贴自我的设备", notify); err != nil {
		requestLogger.WithError(err).WithField("notify", notify).Warn("failed to send notification")
	}

	w.WriteHeader(http.StatusOK)
	return
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	requestLogger := log.WithFields(log.Fields{"request_id": rand.Int(), "user_ip": r.RemoteAddr})
	requestLogger.Info("404 not found")
	w.WriteHeader(http.StatusNotFound)
}

func getCurrentPath() string {
	dir, _ := os.Getwd()
	return dir
}

func getTempFilePath(filename string) string {
	if !filepath.IsAbs(config.GetTempDir()) {
		p, err := filepath.Abs(config.GetTempDir())
		if err != nil {
			return filepath.Join(getCurrentPath(), config.GetTempDir(), filename)
		}
		return filepath.Join(p, filename)
	}
	return filepath.Join(config.GetTempDir(), filename)
}

func setLastFilename(filename string) {
	path := getTempFilePath("_filename.txt")
	_ = ioutil.WriteFile(path, []byte(filename), os.ModePerm)
}

func cleanTempFile() {
	tempDir := getTempFilePath("")
	if a, err := os.Stat(tempDir); err != nil || !a.IsDir() {
		_ = os.Mkdir(tempDir, os.ModePerm)
	}

	path := getTempFilePath("_filename.txt")
	if isExistFile(path) {
		filename, err := ioutil.ReadFile(path)
		if err != nil || len(filename) == 0 {
			return
		}

		delPath := getTempFilePath(string(filename))
		if isExistFile(delPath) {
			_ = os.Remove(delPath)
		}

		setLastFilename("")
	}
}
