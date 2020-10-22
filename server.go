package main

import (
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/YanxinTang/clipboard-online/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	typeText  = "text"
	typeFile  = "file"
	typeMedia = "media"
)

func setupRoute(engin *gin.Engine) {
	// engin.GET("/", getHandler)
	engin.POST("/", setHandler)
	engin.NoRoute(notFoundHandler)
}

// func getHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	requestLogger := log.WithFields(log.Fields{"request_id": rand.Int(), "user_ip": r.RemoteAddr})

// 	contentType, err := utils.Clipboard().ContentType()
// 	if err != nil {
// 		requestLogger.WithError(err).Info("failed to get content type of clipboard")
// 		return
// 	}
// 	requestLogger.WithField("content type", contentType).Info("get content type of clipboard")

// 	if contentType == typeText {
// 		str, err := walk.Clipboard().Text()
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			_, _ = io.WriteString(w, "")
// 			requestLogger.WithError(err).Warn("failed to get clipboard")
// 			return
// 		}
// 		writeJSON(w, H{
// 			"type": "text",
// 			"data": str,
// 		})
// 		requestLogger.Info("get clipboard text")
// 		return
// 	}

// 	if contentType == typeFile {
// 		filenames, err := utils.Clipboard().Files()
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		type ResponseFile struct {
// 			Name    string `json:"name"`
// 			Content string `json:"content"`
// 		}

// 		type ResponseFiles []ResponseFile

// 		responseFiles := make([]ResponseFile, 0, len(filenames))
// 		for _, path := range filenames {
// 			base64, err := readBase64FromFile(path)
// 			if err != nil {
// 				log.WithError(err).WithField("filepath", path).Warning("read base64 from file failed")
// 				continue
// 			}
// 			responseFiles = append(responseFiles, ResponseFile{filepath.Base(path), base64})
// 		}

// 		writeJSON(w, H{
// 			"type": "file",
// 			"data": responseFiles,
// 		})
// 		requestLogger.Info("get clipboard files")
// 		return
// 	}
// 	w.WriteHeader(http.StatusBadRequest)
// }

func readBase64FromFile(path string) (string, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

type Body struct {
	Type               string `json:"type"`
	NamesString        string `json:"names"`
	EncodedFilesString string `json:"files"`
}

type ByteFile struct {
	Name  string
	Bytes []byte
}

func (b *Body) Names() []string {
	return strings.Split(b.NamesString, "\n")
}

func (b *Body) ByteFiles() []ByteFile {
	encodedFiles := strings.Split(b.EncodedFilesString, "\n")
	byteFiles := make([]ByteFile, 0, len(encodedFiles))
	names := b.Names()
	for i, encodedFile := range encodedFiles {
		fileBytes, err := base64.StdEncoding.DecodeString(encodedFile)
		if err != nil {
			log.WithError(err).Warn("failed to decode file base64")
			continue
		}
		byteFiles = append(byteFiles, ByteFile{names[i], fileBytes})
	}
	return byteFiles
}

func setHandler(c *gin.Context) {
	requestLogger := log.WithFields(log.Fields{"request_id": rand.Int(), "user_ip": c.Request.RemoteAddr})
	cleanTempFile()

	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		requestLogger.WithError(err).Warn("failed to bind body")
		c.Status(http.StatusBadRequest)
		return
	}

	if body.Type == typeText {
		c.Status(http.StatusOK)
		return
	}

	if body.Type == typeFile {
		byteFiles := body.ByteFiles()
		paths := make([]string, 0, len(byteFiles))
		for _, byteFile := range byteFiles {
			path := getTempFilePath(string(byteFile.Name))
			if err := ioutil.WriteFile(path, byteFile.Bytes, 0644); err != nil {
				requestLogger.WithError(err).WithField("path", path).Warn("failed to create file")
				continue
			}
			paths = append(paths, path)
		}

		if err := utils.Clipboard().SetFiles(paths); err != nil {
			requestLogger.WithError(err).Warn("failed to set clipboard")
			c.Status(http.StatusBadRequest)
			return
		}

		requestLogger.WithField("paths", paths).Info("set clipboard file")
		c.Status(http.StatusOK)
		return
	}

	// rd := bufio.NewReader(r.Body)

	// var (
	// 	version     string
	// 	contentType string
	// 	notify      string
	// 	filename    string
	// )

	// q := r.URL.Query()
	// version = q.Get("version")
	// contentType = q.Get("type")
	// filename = q.Get("filename")

	// if len(filename) == 0 {
	// 	filename = utils.RandStringBytes(16)
	// }

	// cleanTempFile()

	// if version == "" || contentType == typeText {
	// 	text, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		requestLogger.WithError(err).Warn("failed to read request body")
	// 		return
	// 	}

	// 	if len(text) == 0 {
	// 		notify = "粘贴内容为空"
	// 	} else {
	// 		notify = string(text)
	// 	}

	// 	if err := walk.Clipboard().SetText(string(text)); err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		requestLogger.WithError(err).Warn("failed to set clipboard")
	// 		return
	// 	}

	// 	requestLogger.WithField("text", string(text)).Info("set clipboard text")
	// } else if contentType == typeFile || contentType == typeMedia {
	// 	path := getTempFilePath(filename)
	// 	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		requestLogger.WithError(err).Warn("failed create temporary file")
	// 		return
	// 	}
	// 	defer file.Close()

	// 	_, _ = io.Copy(file, rd)

	// 	if contentType == typeMedia {
	// 		notify = "[图片媒体] 已复制到剪贴板"
	// 	} else {
	// 		notify = "[文件] 已复制到剪贴板"
	// 	}

	// 	setLastFilename(filename)
	// 	paths := []string{path, `D:\Projects\golang\clipboard-online\images\clipboard-icon.png`}
	// 	if err := utils.Clipboard().SetFiles(paths); err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		requestLogger.WithError(err).Warn("failed to set clipboard")
	// 		return
	// 	}

	// 	requestLogger.WithField("file", path).Info("set clipboard file")
	// } else {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	requestLogger.Warn("unsupported content type")
	// 	return
	// }

	// if err := app.ni.ShowInfo("粘贴自我的设备", notify); err != nil {
	// 	requestLogger.WithError(err).WithField("notify", notify).Warn("failed to send notification")
	// }

	// w.WriteHeader(http.StatusOK)
	// return
}

func notFoundHandler(c *gin.Context) {
	requestLogger := log.WithFields(log.Fields{"request_id": rand.Int(), "user_ip": c.Request.RemoteAddr})
	requestLogger.Info("404 not found")
	c.Status(http.StatusNotFound)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
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
