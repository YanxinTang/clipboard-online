package main

const ConfigFilePath = "config.json"

type Config struct {
	Port string `json:"port"`
}

var DefaultConfig = Config{
	Port: "8086",
}
