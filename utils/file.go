package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func IsExistFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDirectory(path string) error {
	if a, err := os.Stat(path); err != nil || !a.IsDir() {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func AppendOrderToFilename(path string) string {
	dirPath := filepath.Dir(path)
	basename := filepath.Base(path) // filename with ext
	ext := filepath.Ext(basename)
	name := strings.TrimSuffix(basename, ext)

	newName := fmt.Sprintf("%s(%d)%s", name, 1, ext)

	reg := regexp.MustCompile(`^(.*)\((\d+)\)$`)
	matchResult := reg.FindSubmatch([]byte(name))

	if len(matchResult) == 3 {
		originName := string(matchResult[1])
		lastOrder, err := strconv.Atoi(string(matchResult[2]))

		if err == nil {
			newName = fmt.Sprintf("%s(%d)%s", originName, lastOrder+1, ext)
		}
	}
	return filepath.Join(dirPath, newName)
}

func LatestFilename(path string) string {
	if IsExistFile(path) {
		latest := AppendOrderToFilename(path)
		return LatestFilename(latest)
	}
	return path
}
