package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	WorkDirs []string `json:"workDirs"`
}

func NewConfig() Config {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}
