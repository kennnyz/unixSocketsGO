package configs

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ListenAddress     string `json:"listenAddress"`
	MsgChanBufferSize int    `json:"msgChanBufferSize"`
}

func ReadConfig(configPath string) Config {
	config := Config{}
	file, err := os.Open(configPath)
	if err != nil {
		log.Println("error reading config", err)
		return Config{}
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Println("error reading config")
		return Config{}
	}
	return config
}
