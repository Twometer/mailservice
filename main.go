package main

import "log"

func main() {
	config, err := LoadConfig("./config.example.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	log.Println("Loaded config", config)
}
