package main

const ConfigFile = "config.json"
const LogFile = "log.txt"

type Config struct {
	Port string `json:"port"`
}

var DefaultConfig = Config{
	Port: "8086",
}
