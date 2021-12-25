package main

import (
	"encoding/json"
	"os"
)

type SiteConfig struct {
	ReceiverAddress string
	CorsOrigin      string
	Fields          []string
}

type ServiceConfig struct {
	SenderAddress  string
	SenderPassword string
	ServerName     string
	ServerPort     uint16
	Sites          map[string]SiteConfig
}

func LoadConfig(path string) (ServiceConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return ServiceConfig{}, err
	}

	var config ServiceConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return ServiceConfig{}, err
	}
	return config, nil
}
