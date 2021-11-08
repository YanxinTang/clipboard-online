package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/YanxinTang/clipboard-online/utils"
	"github.com/sirupsen/logrus"
)

const ConfigFile = "config.json"
const LogFile = "log.txt"

// Config represents configuration for applicaton
type Config struct {
	Port                  string       `json:"port"`
	Authkey               string       `json:"authkey"`
	AuthkeyExpiredTimeout int64        `json:"authkeyExpiredTimeout"`
	LogLevel              logrus.Level `json:"logLevel"`
	TempDir               string       `json:"tempDir"`
	ReserveHistory        bool         `json:"reserveHistory"`
	Notify                ConfigNotify `json:"notify"`
}

type ConfigNotify struct {
	Copy  bool `json:"copy"`
	Paste bool `json:"paste"`
}

// DefaultConfig is a default configuration for application
var DefaultConfig = Config{
	Port:                  "8086",
	Authkey:               "",
	AuthkeyExpiredTimeout: 30,
	LogLevel:              logrus.WarnLevel,
	TempDir:               "./temp",
	ReserveHistory:        false,
	Notify: ConfigNotify{
		Copy:  false,
		Paste: false,
	},
}

func loadConfig(path string) (*Config, error) {
	if utils.IsExistFile(path) {
		return loadConfigFromFile(path)
	}
	if err := createConfigFile(path); err != nil {
		return nil, err
	}
	return &DefaultConfig, nil
}

func loadConfigFromFile(path string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(configBytes, &DefaultConfig); err != nil {
		return nil, err
	}
	return &DefaultConfig, nil
}

func createConfigFile(path string) error {
	defaultConfigJSON, err := json.MarshalIndent(DefaultConfig, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, []byte(defaultConfigJSON), 0744); err != nil {
		return err
	}
	return nil
}
