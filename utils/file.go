package utils

import "os"

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
