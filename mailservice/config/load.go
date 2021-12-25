package config

import (
	"encoding/json"
	"os"
)

func Read(path string) (ServiceConfig, error) {
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
